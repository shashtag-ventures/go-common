package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode >= 400 {
		rw.body.Write(b)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func getRealIP(r *http.Request) string {
	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

func scrubPayload(payload []byte) string {
	if len(payload) == 0 {
		return ""
	}
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return string(payload)
	}

	sensitiveKeys := []string{"password", "token", "secret", "access_token", "refresh_token"}
	var scrub func(m map[string]interface{})
	scrub = func(m map[string]interface{}) {
		for k, v := range m {
			for _, sk := range sensitiveKeys {
				if strings.EqualFold(k, sk) {
					m[k] = "[MASKED]"
				}
			}
			if child, ok := v.(map[string]interface{}); ok {
				scrub(child)
			}
		}
	}
	scrub(data)
	scrubbed, _ := json.Marshal(data)
	return string(scrubbed)
}

func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v1/health" || r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			var reqBody []byte
			if r.Body != nil && r.Method != http.MethodGet {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			rw := &responseWriter{w, http.StatusOK, bytes.NewBuffer(nil), 0}
			next.ServeHTTP(rw, r)

			ctx := r.Context()
			duration := time.Since(start)
			ms := duration.Milliseconds()

			// 1. Resolve State (The Mutable Fix)
			var userID uint
			var extraFields map[string]any
			if state, ok := ctx.Value(LogStateKey).(*LogState); ok {
				userID = state.UserID
				extraFields = state.Fields
			}

			// 2. Logic Escalation
			level := slog.LevelInfo
			if rw.statusCode >= 500 {
				level = slog.LevelError
			} else if rw.statusCode >= 400 || ms > 500 {
				level = slog.LevelWarn
			}

			// 3. Metadata extraction
			requestID, _ := ctx.Value(RequestIDKey).(string)
			var traceID string
			if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}

			// 4. Construct Final Grouped Log
			attrs := []any{
				slog.Group("http",
					slog.String("method", r.Method),
					slog.String("url", r.URL.String()),
					slog.Int("status", rw.statusCode),
					slog.Int64("duration_ms", ms),
					slog.Int("size_bytes", rw.size),
					slog.String("ua", r.UserAgent()),
					slog.String("referer", r.Referer()),
				),
				slog.Group("user",
					slog.Uint64("id", uint64(userID)),
					slog.String("ip", getRealIP(r)),
				),
				slog.Group("trace",
					slog.String("request_id", requestID),
					slog.String("trace_id", traceID),
				),
			}

			// Add enriched fields from handlers
			if len(extraFields) > 0 {
				appFields := []any{}
				for k, v := range extraFields {
					appFields = append(appFields, slog.Any(k, v))
				}
				attrs = append(attrs, slog.Group("app", appFields...))
			}

			if len(reqBody) > 0 && len(reqBody) < 4096 {
				attrs = append(attrs, slog.String("payload", scrubPayload(reqBody)))
			}

			if rw.statusCode >= 400 && rw.body.Len() > 0 {
				attrs = append(attrs, slog.String("error", rw.body.String()))
			}

			slog.Log(ctx, level, "HTTP Request", attrs...)
		})
	}
}
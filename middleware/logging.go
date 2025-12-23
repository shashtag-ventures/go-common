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
			requestID, _ := ctx.Value(RequestIDKey).(string)
			
			var traceID string
			if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}

			var userID uint
			if user, ok := ctx.Value(UserContextKey).(*AuthenticatedUser); ok {
				userID = user.ID
			}

			duration := time.Since(start)

			// LOG LEVEL ESCALATION: Auto-detect importance
			level := slog.LevelInfo
			if rw.statusCode >= 500 {
				level = slog.LevelError
			} else if rw.statusCode >= 400 || duration > 500*time.Millisecond {
				level = slog.LevelWarn
			}

			attrs := []any{
				"method",      r.Method,
				"url",         r.URL.String(),
				"status",      rw.statusCode,
				"duration_ms", duration.Milliseconds(),
				"size_bytes",  rw.size, // Bandwidth tracking
				"user_id",     userID,
				"request_id",  requestID,
				"trace_id",    traceID,
				"ip",          r.RemoteAddr,
				"user_agent",  r.UserAgent(),
				"referer",     r.Referer(), // RESTORED
			}

			if len(reqBody) > 0 && len(reqBody) < 4096 {
				attrs = append(attrs, "payload", scrubPayload(reqBody))
			}

			if rw.statusCode >= 400 && rw.body.Len() > 0 {
				attrs = append(attrs, "error_detail", rw.body.String())
			}

			slog.Log(ctx, level, "HTTP Request", attrs...)
		})
	}
}
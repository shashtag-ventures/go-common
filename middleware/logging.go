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
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode >= 400 {
		rw.body.Write(b)
	}
	return rw.ResponseWriter.Write(b)
}

// scrubPayload parses JSON and masks sensitive keys within the payload string.
func scrubPayload(payload []byte) string {
	if len(payload) == 0 {
		return ""
	}
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return string(payload) // Return as-is if not valid JSON
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
			// Skip logging for successful health checks to reduce noise
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

			rw := &responseWriter{w, http.StatusOK, bytes.NewBuffer(nil)}
			next.ServeHTTP(rw, r)

			requestID, _ := r.Context().Value(RequestIDKey).(string)
			var traceID string
			if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}

			var userID uint
			if user, ok := r.Context().Value(UserContextKey).(*AuthenticatedUser); ok {
				userID = user.ID
			}

			attrs := []any{
				"method",      r.Method,
				"url",         r.URL.String(),
				"status",      rw.statusCode,
				"duration",    time.Since(start).String(),
				"user_id",     userID,
				"request_id",  requestID,
				"trace_id",    traceID,
				"ip",          r.RemoteAddr,
				"user_agent",  r.UserAgent(),
				"referer",     r.Referer(),
			}

			if len(reqBody) > 0 && len(reqBody) < 4096 {
				attrs = append(attrs, "request_payload", scrubPayload(reqBody))
			}

			if rw.statusCode >= 400 && rw.body.Len() > 0 {
				attrs = append(attrs, "response_error", rw.body.String())
			}

			slog.Info("HTTP Request", attrs...)
		})
	}
}

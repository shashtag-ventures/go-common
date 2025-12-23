package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
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
		rw.body.Write(b) // Capture error responses
	}
	return rw.ResponseWriter.Write(b)
}

func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 1. Capture Request Body safely
			var reqBody []byte
			if r.Body != nil && r.Method != http.MethodGet {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset body for handler
			}

			rw := &responseWriter{w, http.StatusOK, bytes.NewBuffer(nil)}
			requestID, _ := r.Context().Value(RequestIDKey).(string)

			var traceID string
			if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}

			next.ServeHTTP(rw, r)

			var userID uint
			if user, ok := r.Context().Value(UserContextKey).(*AuthenticatedUser); ok {
				userID = user.ID
			}

			// Prepare log attributes
			attrs := []any{
				"method", r.Method,
				"url", r.URL.String(),
				"status", rw.statusCode,
				"duration", time.Since(start).String(),
				"user_id", userID,
				"request_id", requestID,
				"trace_id", traceID,
				"ip", r.RemoteAddr,
			}

			// Only log body if it exists and isn't too large
			if len(reqBody) > 0 && len(reqBody) < 2048 {
				attrs = append(attrs, "request_payload", string(reqBody))
			}

			// Log response error message if it failed
			if rw.statusCode >= 400 && rw.body.Len() > 0 {
				attrs = append(attrs, "response_error", rw.body.String())
			}

			slog.Info("HTTP Request", attrs...)
		})
	}
}
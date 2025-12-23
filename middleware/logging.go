package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that intercepts the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// RequestLogger returns a middleware that logs the full details of every request.
func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Wrap the response writer to capture the status code
			rw := &responseWriter{w, http.StatusOK}

			// 1. Extract request ID (Your custom ID)
			requestID, _ := r.Context().Value(RequestIDKey).(string)

			// 2. Extract OpenTelemetry Trace ID (The industry standard)
			var traceID string
			if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}

			// Process the request
			next.ServeHTTP(rw, r)

			// 3. Extract user ID if auth has run
			var userID uint
			if user, ok := r.Context().Value(UserContextKey).(*AuthenticatedUser); ok {
				userID = user.ID
			}

			// Log the completed request with high-fidelity metadata
			slog.Info("HTTP Request",
				"method", r.Method,
				"url", r.URL.String(),
				"status", rw.statusCode,
				"duration", time.Since(start).String(),
				"user_id", userID,
				"request_id", requestID,
				"trace_id", traceID, // Links logs to OTel traces
				"ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"referer", r.Referer(),
			)
		})
	}
}

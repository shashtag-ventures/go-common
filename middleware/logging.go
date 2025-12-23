package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that intercepts the status code and captures error bodies.
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

// RequestLogger returns a middleware that logs the full details of every request.
func RequestLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 1. Capture Request Body safely for non-GET requests
			var reqBody []byte
			if r.Body != nil && r.Method != http.MethodGet {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			rw := &responseWriter{w, http.StatusOK, bytes.NewBuffer(nil)}

			// Process the request
			next.ServeHTTP(rw, r)

			// 2. Metadata Extraction (After next.ServeHTTP so context is fully populated if possible)
			requestID, _ := r.Context().Value(RequestIDKey).(string)
			
			var traceID string
			if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}

			var userID uint
			if user, ok := r.Context().Value(UserContextKey).(*AuthenticatedUser); ok {
				userID = user.ID
			}

			// 3. Prepare log attributes
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

			// Add payload if it exists and is reasonable in size
			if len(reqBody) > 0 && len(reqBody) < 2048 {
				attrs = append(attrs, "request_payload", string(reqBody))
			}

			// Add error response if applicable
			if rw.statusCode >= 400 && rw.body.Len() > 0 {
				attrs = append(attrs, "response_error", rw.body.String())
			}

			slog.Info("HTTP Request", attrs...)
		})
	}
}

package middleware

import (
	"log/slog"
	"net/http"
	"time"
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

			// Extract request ID if present in context (from RequestID middleware)
			requestID, _ := r.Context().Value(RequestIDKey).(string)

			// Process the request
			next.ServeHTTP(rw, r)

			// Post-processing: try to extract user ID if auth middleware has run
			var userID uint
			if user, ok := r.Context().Value(UserContextKey).(*AuthenticatedUser); ok {
				userID = user.ID
			}

			// Log the completed request with full original URL and duration
			slog.Info("HTTP Request",
				"method", r.Method,
				"url", r.URL.String(),
				"status", rw.statusCode,
				"duration", time.Since(start).String(),
				"user_id", userID,
				"request_id", requestID,
				"ip", r.RemoteAddr,
			)
		})
	}
}
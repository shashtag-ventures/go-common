package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// RequestIDMiddleware generates a unique request ID for each incoming request
// and adds it to the request's context and the slog logger.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		// Add the request ID to the request context.
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		// Add the request ID to the slog logger for this request.
		// This creates a new logger with the request ID as an attribute.
		logger := slog.Default().With(slog.String("requestID", requestID))
		ctx = context.WithValue(ctx, LoggerContextKey, logger) // Store the new logger in context
		r = r.WithContext(ctx)

		// Set the X-Request-ID header in the response.
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r)
	})
}

// GetLoggerFromContext retrieves the slog.Logger from the request context.
// If no logger is found, it returns the default logger.
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(LoggerContextKey).(*slog.Logger)
	if !ok || logger == nil {
		return slog.Default()
	}
	return logger
}

// GetRequestIDFromContext retrieves the request ID from the request context.
func GetRequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

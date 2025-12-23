package middleware

import (
	"context"
	"log/slog"
)

// CtxKey is a custom type for context keys to avoid collisions.
type CtxKey string

const (
	// RequestIDKey is the key used to store the request ID in the context.
	RequestIDKey CtxKey = "requestID"
	// LoggerContextKey is the key used to store the slog.Logger in the context.
	LoggerContextKey CtxKey = "logger"
	// UserContextKey is the key used to store the authenticated user in the context.
	UserContextKey CtxKey = "user"
)

// GetLoggerFromContext retrieves the logger from the context.
// If no logger is found, it returns the default logger.
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// GetRequestIDFromContext retrieves the request ID from the context.
func GetRequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
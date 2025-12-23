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
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// EnrichLogger adds attributes to the logger stored in the context and returns a new context.
// Use this inside handlers to add metadata (like project_id) that should appear in all subsequent logs.
func EnrichLogger(ctx context.Context, attrs ...slog.Attr) context.Context {
	logger := GetLoggerFromContext(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}
	return context.WithValue(ctx, LoggerContextKey, logger)
}

// GetRequestIDFromContext retrieves the request ID from the context.
func GetRequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

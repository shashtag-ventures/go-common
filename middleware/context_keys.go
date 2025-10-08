package middleware

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

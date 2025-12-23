package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// RequestIDMiddleware is a middleware that injects a unique request ID into the context of each request.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", requestID)
		
		// 1. Store the raw ID in context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		
		// 2. Create a logger that ALWAYS includes the request_id
		// This means every slog call in the handler automatically has the ID
		logger := slog.Default().With("request_id", requestID)
		ctx = context.WithValue(ctx, LoggerContextKey, logger)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
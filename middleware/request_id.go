package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", requestID)
		
		// 1. Initialize the mutable Log State
		state := &LogState{
			Fields: make(map[string]any),
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		ctx = context.WithValue(ctx, LogStateKey, state) // NEW: Pass state pointer
		
		logger := slog.Default().With("request_id", requestID)
		ctx = context.WithValue(ctx, LoggerContextKey, logger)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

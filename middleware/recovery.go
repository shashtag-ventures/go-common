package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery is a middleware that recovers from panics, logs the panic (and a stack trace), 
// and returns a HTTP 500 (Internal Server Error) status.
func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Extract request ID for context
					requestID, _ := r.Context().Value(RequestIDKey).(string)
					
					slog.Error("PANIC RECOVERED",
						"error", err,
						"request_id", requestID,
						"stack", string(debug.Stack()),
					)
					
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"net/http"
)

// TrailingSlashMiddleware removes a trailing slash from the request URL path if present (unless it's the root path).
// This helps in standardizing URLs and avoiding duplicate content issues.
func TrailingSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is not the root ("/") and ends with a slash.
		if len(r.URL.Path) > 1 && r.URL.Path[len(r.URL.Path)-1] == '/' {
			// Remove the trailing slash.
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-1]
		}
		// Call the next handler in the chain with the potentially modified request.
		next.ServeHTTP(w, r)
	})
}

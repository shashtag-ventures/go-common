package middleware

import (
	"crypto/sha256"
	"net/http"

	"github.com/gorilla/csrf"
)

// CSRFConfig holds the configuration for the CSRF middleware.
type CSRFConfig struct {
	Enabled bool
	Secret  string // Secret key (will be hashed to 32 bytes)
	Secure  bool   // Use true for production (HTTPS)
	Domain  string // Cookie domain
}

// CSRFMiddleware wraps gorilla/csrf to provide CSRF protection.
func CSRFMiddleware(cfg CSRFConfig) func(http.Handler) http.Handler {
	if !cfg.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	// gorilla/csrf requires exactly 32 bytes. We hash the user-provided secret
	// to ensure it's always the correct length regardless of the env variable value.
	key := sha256.Sum256([]byte(cfg.Secret))

	return csrf.Protect(
		key[:],
		csrf.Secure(cfg.Secure),
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Path("/"),
		csrf.Domain(cfg.Domain),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("csrf_token"),
	)
}

// GetCSRFToken returns the CSRF token for the current request.
// This can be used in a handler to send the token to the client.
func GetCSRFToken(r *http.Request) string {
	return csrf.Token(r)
}

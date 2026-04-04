package middleware

import (
	"crypto/sha256"
	"net/http"
	"net/url"

	"github.com/gorilla/csrf"
)

// CSRFConfig holds the configuration for the CSRF middleware.
type CSRFConfig struct {
	Enabled        bool
	Secret         string   // Secret key (will be hashed to 32 bytes)
	Secure         bool     // Use true for production (HTTPS)
	Domain         string   // Cookie domain
	TrustedOrigins []string // Origins allowed to send state-changing requests (e.g., "https://www.example.com")
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

	opts := []csrf.Option{
		csrf.Secure(cfg.Secure),
		csrf.HttpOnly(true),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.Path("/"),
		csrf.Domain(cfg.Domain),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("csrf_token"),
	}

	if len(cfg.TrustedOrigins) > 0 {
		// gorilla/csrf compares TrustedOrigins against referer.Host (bare hostname),
		// so we must strip the scheme from full URLs like "https://www.example.com".
		hosts := make([]string, 0, len(cfg.TrustedOrigins))
		for _, origin := range cfg.TrustedOrigins {
			if parsed, err := url.Parse(origin); err == nil && parsed.Host != "" {
				hosts = append(hosts, parsed.Host)
			} else {
				// Already a bare host or unparseable — pass through as-is.
				hosts = append(hosts, origin)
			}
		}
		opts = append(opts, csrf.TrustedOrigins(hosts))
	}

	return csrf.Protect(key[:], opts...)
}

// GetCSRFToken returns the CSRF token for the current request.
// This can be used in a handler to send the token to the client.
func GetCSRFToken(r *http.Request) string {
	return csrf.Token(r)
}

package middleware

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/shashtag-ventures/go-common/jsonResponse"
)

// RateLimitConfig defines settings for API rate limiting.
type RateLimitConfig struct {
	Enabled bool
	Limit   int
	Window  int // Window in seconds
}

// client represents a client for rate limiting, tracking request count and last request time.
type client struct {
	lastRequest  time.Time
	requestCount int
}

// RateLimitMiddleware provides a basic IP-based rate limiting.
// DEPRECATED: Use RateLimitMiddlewareWithContext instead to avoid goroutine leaks.
func RateLimitMiddleware(cfg RateLimitConfig) func(next http.Handler) http.Handler {
	return RateLimitMiddlewareWithContext(context.Background(), cfg)
}

// RateLimitMiddlewareWithContext provides a basic IP-based rate limiting with context for cleanup.
func RateLimitMiddlewareWithContext(ctx context.Context, cfg RateLimitConfig) func(next http.Handler) http.Handler {
	mu := &sync.Mutex{}
	clients := make(map[string]*client)

	if cfg.Enabled {
		// Goroutine to periodically clean up old client entries from the map.
		go func() {
			ticker := time.NewTicker(time.Duration(cfg.Window) * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					mu.Lock()
					for ip, c := range clients {
						if time.Since(c.lastRequest) > time.Duration(cfg.Window)*time.Second {
							delete(clients, ip)
						}
					}
					mu.Unlock()
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			ip := r.RemoteAddr

			mu.Lock()
			c, found := clients[ip]
			if !found {
				c = &client{lastRequest: time.Now(), requestCount: 0}
				clients[ip] = c
			}

			if time.Since(c.lastRequest) > time.Duration(cfg.Window)*time.Second {
				c.requestCount = 0
				c.lastRequest = time.Now()
			}

			if c.requestCount >= cfg.Limit {
				mu.Unlock()
				jsonResponse.SendErrorResponse(w, errors.New(http.StatusText(http.StatusTooManyRequests)), http.StatusTooManyRequests)
				return
			}

			c.requestCount++
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

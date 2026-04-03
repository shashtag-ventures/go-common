package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/shashtag-ventures/go-common/jsonResponse"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	limiterredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RateLimitStore defines the type of storage to use for rate limiting counters.
type RateLimitStore string

const (
	RateLimitStoreMemory RateLimitStore = "memory"
	RateLimitStoreRedis  RateLimitStore = "redis"
)

// RateLimitConfig defines settings for API rate limiting.
type RateLimitConfig struct {
	Enabled     bool
	Limit       int
	Window      int             // Window in seconds
	StoreType   RateLimitStore  // "memory" or "redis"
	RedisClient limiterredis.Client // Required for "redis" store
	KeyPrefix   string          // Optional prefix for rate limit keys
}

// RateLimitMiddleware provides a basic IP-based rate limiting using ulule/limiter.
// It supports both in-memory (for simple cases) and Redis (for distributed systems).
func RateLimitMiddleware(cfg RateLimitConfig) func(next http.Handler) http.Handler {
	if !cfg.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	// Define the rate
	rate := limiter.Rate{
		Period: time.Duration(cfg.Window) * time.Second,
		Limit:  int64(cfg.Limit),
	}

	// Initialize the store
	var store limiter.Store
	var err error

	switch cfg.StoreType {
	case RateLimitStoreRedis:
		if cfg.RedisClient == nil {
			panic("RateLimitConfig: RedisClient is required when StoreType is 'redis'")
		}
		store, err = limiterredis.NewStoreWithOptions(cfg.RedisClient, limiter.StoreOptions{
			Prefix: cfg.KeyPrefix,
		})
	default:
		// Default to memory
		store = memory.NewStore()
	}

	if err != nil {
		// If store initialization fails, we panic in production to avoid unprotected endpoints,
		// or we could fall back to memory. For now, panic for visibility.
		panic(fmt.Sprintf("RateLimitConfig: failed to initialize store: %v", err))
	}

	// Create the limiter instance
	instance := limiter.New(store, rate)

	// Create the stdlib middleware
	// By default, it uses the IP address as the key.
	mw := stdlib.NewMiddleware(instance)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We wrap the stdlib middleware handler
			mw.Handler(next).ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddlewareWithContext is kept for backward compatibility but context is ignored 
// as ulule/limiter handles cleanup internally (memory store) or via TTL (redis).
func RateLimitMiddlewareWithContext(ctx context.Context, cfg RateLimitConfig) func(next http.Handler) http.Handler {
	return RateLimitMiddleware(cfg)
}

// CustomRateLimitMiddleware allows for more advanced use cases, like different keys or rates.
func CustomRateLimitMiddleware(instance *limiter.Limiter, keyExtractor func(r *http.Request) string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyExtractor(r)
			ctx := r.Context()

			limitContext, err := instance.Get(ctx, key)
			if err != nil {
				// On error, we allow the request but log the failure
				// In a stricter system, you might want to block
				next.ServeHTTP(w, r)
				return
			}

			// Set standard rate limit headers
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limitContext.Limit))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limitContext.Remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", limitContext.Reset))

			if limitContext.Reached {
				jsonResponse.SendErrorResponse(w, fmt.Errorf("rate limit reached"), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

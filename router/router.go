package router

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shashtag-ventures/go-common/middleware"
)

// Config holds the configuration for the router.
type Config struct {
	ApiVersion      string
	Cors            middleware.CorsConfig
	CSRF            middleware.CSRFConfig
	RateLimit       middleware.RateLimitConfig
}

// Router is a wrapper around http.ServeMux that supports middleware via .Use()
type Router struct {
	*http.ServeMux
	middlewares []func(http.Handler) http.Handler
	config      Config
	once        sync.Once
	fullHandler http.Handler
}

// Use adds middleware(s) to the router.
func (r *Router) Use(m ...func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, m...)
}

// buildHandler wraps the base ServeMux with all configured and custom middlewares.
func (r *Router) buildHandler() http.Handler {
	var handler http.Handler = r.ServeMux

	// Apply custom middlewares added via .Use()
	// Reverse order for standard middleware wrap logic (last applied is first executed)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	// PRODUCTION-ORDER STACK:
	// Logic: Last applied is FIRST executed.

	// 5. Security & Redirects (Inner-most system layers)
	handler = middleware.CSRFMiddleware(r.config.CSRF)(handler)
	handler = middleware.RateLimitMiddleware(r.config.RateLimit)(handler)
	handler = middleware.TrailingSlashMiddleware(handler)
	handler = middleware.CorsMiddleware(r.config.Cors, handler)

	// 4. Observability (Capture metrics for the secured request)
	handler = middleware.MetricsMiddleware(handler)

	// 3. Panic Recovery (Protect monitoring layers from handler crashes)
	handler = middleware.Recovery()(handler)

	// 2. Global Request Logger (Must run after ID is set)
	handler = middleware.RequestLogger()(handler)

	// 1. Assign Request ID (Outer-most layer - must run first)
	handler = middleware.RequestIDMiddleware(handler)

	return handler
}

// ServeHTTP implements http.Handler and lazily builds the middleware chain.
func (r *Router) ServeHTTP(w http.ResponseWriter, r2 *http.Request) {
	r.once.Do(func() {
		r.fullHandler = r.buildHandler()
	})

	apiPath := "/api/" + r.config.ApiVersion
	http.StripPrefix(apiPath, r.fullHandler).ServeHTTP(w, r2)
}

// New creates and configures the main and API routers with standard middleware.
func New(cfg Config) (*http.ServeMux, *Router) {
	mainRouter := http.NewServeMux()
	apiMux := http.NewServeMux()

	apiRouter := &Router{
		ServeMux: apiMux,
		config:   cfg,
	}

	// Handle all API requests through the apiRouter, which manages the middleware chain.
	mainRouter.Handle("/api/"+cfg.ApiVersion+"/", apiRouter)

	mainRouter.Handle("/metrics", promhttp.Handler())

	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is healthy"))
	})

	return mainRouter, apiRouter
}
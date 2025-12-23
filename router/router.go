package router

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shashtag-ventures/go-common/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Config holds the configuration for the router.
type Config struct {
	ApiVersion      string
	Cors            middleware.CorsConfig
	RateLimit       middleware.RateLimitConfig
	OtelServiceName string
}

// Router is a wrapper around http.ServeMux that supports middleware via .Use()
type Router struct {
	*http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

// Use adds middleware(s) to the router.
func (r *Router) Use(m ...func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, m...)
}

// New creates and configures the main and API routers with standard middleware.
func New(cfg Config) (*http.ServeMux, *Router) {
	mainRouter := http.NewServeMux()
	apiMux := http.NewServeMux()

	apiRouter := &Router{
		ServeMux: apiMux,
	}

	// Wrapper to apply all middlewares
	mainRouter.HandleFunc("/api/"+cfg.ApiVersion+"/", func(w http.ResponseWriter, r *http.Request) {
		var handler http.Handler = apiRouter.ServeMux

		// Apply custom middlewares added via .Use() (Deepest level)
		for i := len(apiRouter.middlewares) - 1; i >= 0; i-- {
			handler = apiRouter.middlewares[i](handler)
		}

		// Apply standard system middlewares in the "Perfect Order"
		
		// 1. Assign Request ID first (Outermost)
		handler = middleware.RequestIDMiddleware(handler)
		
		// 2. Global Request Logger (must be outside security layers to log blocks)
		handler = middleware.RequestLogger()(handler)
		
		// 3. Panic Recovery (inner protection)
		handler = middleware.Recovery()(handler)
		
		// 4. Observability
		handler = middleware.MetricsMiddleware(handler)
		handler = otelhttp.NewHandler(handler, cfg.OtelServiceName)
		
		// 5. Security & Redirects
		handler = middleware.RateLimitMiddleware(cfg.RateLimit)(handler)
		handler = middleware.TrailingSlashMiddleware(handler)
		handler = middleware.CorsMiddleware(cfg.Cors, handler)

		apiPath := "/api/" + cfg.ApiVersion
		http.StripPrefix(apiPath, handler).ServeHTTP(w, r)
	})

	mainRouter.Handle("/metrics", promhttp.Handler())

	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is healthy"))
	})

	return mainRouter, apiRouter
}

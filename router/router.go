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

		// Apply custom middlewares added via .Use()
		for i := len(apiRouter.middlewares) - 1; i >= 0; i-- {
			handler = apiRouter.middlewares[i](handler)
		}

		// Apply standard system middlewares in CORRECT PRODUCTION ORDER:
		// 1. Recover from panics (Inner protection)
		handler = middleware.Recovery()(handler)
		// 2. Metrics 
		handler = middleware.MetricsMiddleware(handler)
		// 3. OpenTelemetry 
		handler = otelhttp.NewHandler(handler, cfg.OtelServiceName)
		// 4. Global Request Logger (Captures final status after metrics/telemetry)
		handler = middleware.RequestLogger()(handler)
		// 5. Rate Limiting
		handler = middleware.RateLimitMiddleware(cfg.RateLimit)(handler)
		// 6. Request ID (Assign IDs early)
		handler = middleware.RequestIDMiddleware(handler)
		// 7. Utilities (Trailing slash, CORS)
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
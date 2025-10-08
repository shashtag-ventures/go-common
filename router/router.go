package router

import (
	"net/http"
	"strings"

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

// New creates and configures the main and API routers with standard middleware.
func New(cfg Config) (*http.ServeMux, *http.ServeMux) {
	mainRouter := http.NewServeMux()
	apiRouter := http.NewServeMux()

	var apiHandler http.Handler = apiRouter
	apiHandler = middleware.MetricsMiddleware(apiHandler)
	apiHandler = middleware.RateLimitMiddleware(cfg.RateLimit)(apiHandler)
	apiHandler = otelhttp.NewHandler(apiHandler, cfg.OtelServiceName)
	apiHandler = middleware.RequestIDMiddleware(apiHandler)
	apiHandler = middleware.TrailingSlashMiddleware(apiHandler)
	apiHandler = middleware.CorsMiddleware(cfg.Cors, apiHandler)

	apiPath := "/api/" + cfg.ApiVersion + "/"
	mainRouter.Handle(apiPath, http.StripPrefix(strings.TrimSuffix(apiPath, "/"), apiHandler))
	mainRouter.Handle("/metrics", promhttp.Handler())

	apiRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is healthy"))
	})

	return mainRouter, apiRouter
}

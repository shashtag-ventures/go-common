package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CorsConfig defines the configuration for the CORS middleware.
type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// CorsMiddleware creates a CORS middleware based on the provided configuration.
func CorsMiddleware(cfg CorsConfig, next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		AllowCredentials: cfg.AllowCredentials,
	})

	return c.Handler(next)
}
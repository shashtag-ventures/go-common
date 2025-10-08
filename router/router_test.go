package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/shashtag-ventures/go-common/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	// 1. Setup router config
	cfg := router.Config{
		ApiVersion: "v1",
		Cors: middleware.CorsConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
		},
		RateLimit: middleware.RateLimitConfig{
			Enabled: false, // Disable for this test
		},
		OtelServiceName: "test-service",
	}

	// 2. Create the router
	mainRouter, apiRouter := router.New(cfg)
	require.NotNil(t, mainRouter)
	require.NotNil(t, apiRouter)

	// 3. Create a test server
	server := httptest.NewServer(mainRouter)
	defer server.Close()

	t.Run("Health Check Endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "API is healthy", string(body))
	})

	t.Run("Metrics Endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/metrics")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		// Check for a known Prometheus metric
		assert.Contains(t, string(body), "go_goroutines")
	})

	t.Run("RequestID Middleware", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/v1/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		// Check that the middleware added the header
		assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
	})

	t.Run("Trailing Slash Middleware", func(t *testing.T) {
		// Add a temporary handler to test this
		apiRouter.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			// This handler should be hit by a request to "/test" because of the middleware
			w.WriteHeader(http.StatusOK)
		})

		resp, err := http.Get(server.URL + "/api/v1/test/") // Request with trailing slash
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("CORS Preflight", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", server.URL+"/api/v1/health", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	})
}

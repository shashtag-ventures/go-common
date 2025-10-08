package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/stretchr/testify/assert"
)

func TestTrailingSlashMiddleware(t *testing.T) {
	var receivedPath string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	})

	testHandler := middleware.TrailingSlashMiddleware(handler)
	server := httptest.NewServer(testHandler)
	defer server.Close()

	t.Run("Removes trailing slash", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/test/path/")
		assert.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "/test/path", receivedPath)
	})

	t.Run("Does not change root path", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/")
		assert.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "/", receivedPath)
	})

	t.Run("Does not change path with no trailing slash", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/no/slash")
		assert.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "/no/slash", receivedPath)
	})
}

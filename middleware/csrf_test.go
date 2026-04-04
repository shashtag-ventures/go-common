package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCSRFMiddleware(t *testing.T) {
	secret := "32-byte-long-secret-key-for-csrf" // Exactly 32 bytes
	cfg := middleware.CSRFConfig{
		Enabled: true,
		Secret:  secret,
		Secure:  false,
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := middleware.GetCSRFToken(r)
		w.Header().Set("X-CSRF-Token", token)
		w.WriteHeader(http.StatusOK)
	})

	csrfHandler := middleware.CSRFMiddleware(cfg)(nextHandler)

	t.Run("GET Request - Generates Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		csrfHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.NotEmpty(t, rr.Header().Get("X-CSRF-Token"))
		assert.Contains(t, rr.Header().Get("Set-Cookie"), "csrf_token=")
	})

	t.Run("POST Request - Missing Token - Fails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rr := httptest.NewRecorder()

		csrfHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("Disabled CSRF - Succeeds without token", func(t *testing.T) {
		disabledCfg := middleware.CSRFConfig{Enabled: false}
		disabledHandler := middleware.CSRFMiddleware(disabledCfg)(nextHandler)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rr := httptest.NewRecorder()

		disabledHandler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

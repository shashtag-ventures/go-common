package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCSRFMiddleware(t *testing.T) {
	secret := "short-secret" // Less than 32 bytes, hashing should handle it
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

	t.Run("POST Request - Exempt Paths - Succeeds without token", func(t *testing.T) {
		exemptCfg := middleware.CSRFConfig{
			Enabled:     true,
			Secret:      secret,
			ExemptPaths: []string{"/oauth/", "/webhook"},
		}
		exemptHandler := middleware.CSRFMiddleware(exemptCfg)(nextHandler)

		req1 := httptest.NewRequest(http.MethodPost, "/oauth/register", nil)
		rr1 := httptest.NewRecorder()
		exemptHandler.ServeHTTP(rr1, req1)
		assert.Equal(t, http.StatusOK, rr1.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/webhook", nil)
		rr2 := httptest.NewRecorder()
		exemptHandler.ServeHTTP(rr2, req2)
		assert.Equal(t, http.StatusOK, rr2.Code)

		// Non-exempt path should still fail with 403
		req3 := httptest.NewRequest(http.MethodPost, "/projects/create", nil)
		rr3 := httptest.NewRecorder()
		exemptHandler.ServeHTTP(rr3, req3)
		assert.Equal(t, http.StatusForbidden, rr3.Code)
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

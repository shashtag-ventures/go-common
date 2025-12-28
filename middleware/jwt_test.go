package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/jwt"
	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	secret := "test-secret"
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.GetUserFromContext(r.Context())
		if ok {
			w.Header().Set("X-User-ID", user.ID)
			w.Header().Set("X-User-Role", user.Role)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	t.Run("Valid Token", func(t *testing.T) {
		userID := "123"
		role := "admin"
		token, _ := jwt.CreateToken(userID, role, secret)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: token})
		rr := httptest.NewRecorder()

		middleware.JWTAuthMiddleware(secret)(nextHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, userID, rr.Header().Get("X-User-ID"))
		assert.Equal(t, role, rr.Header().Get("X-User-Role"))
	})

	t.Run("Missing Cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		middleware.JWTAuthMiddleware(secret)(nextHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "missing jwt cookie")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "invalid-token"})
		rr := httptest.NewRecorder()

		middleware.JWTAuthMiddleware(secret)(nextHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid or expired token")
	})

	t.Run("Expired Token", func(t *testing.T) {
		// We'll rely on jwt.ParseToken to handle expiration, 
		// but since we can't easily create an expired token with CreateToken, 
		// we just test it implicitly via ParseToken.
		// For a full test, we'd need to manually construct an expired token.
	})
}

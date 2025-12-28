package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestETagMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	})

	etagHandler := ETag(handler)

	t.Run("Generates ETag on first request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		etagHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("ETag"))
		assert.Equal(t, "hello world", w.Body.String())
	})

	t.Run("Returns 304 on matching ETag", func(t *testing.T) {
		// First get the ETag
		req1 := httptest.NewRequest("GET", "/", nil)
		w1 := httptest.NewRecorder()
		etagHandler.ServeHTTP(w1, req1)
		etag := w1.Header().Get("ETag")

		// Second request with If-None-Match
		req2 := httptest.NewRequest("GET", "/", nil)
		req2.Header.Set("If-None-Match", etag)
		w2 := httptest.NewRecorder()
		etagHandler.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusNotModified, w2.Code)
		assert.Empty(t, w2.Body.String())
	})
}

func TestCookieHelpers(t *testing.T) {
	t.Run("SetAuthCookies", func(t *testing.T) {
		w := httptest.NewRecorder()
		SetAuthCookies(w, "test-token", true)

		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 2)

		jwt := cookies[0]
		assert.Equal(t, JWTCookieName, jwt.Name)
		assert.Equal(t, "test-token", jwt.Value)
		assert.True(t, jwt.HttpOnly)
		assert.True(t, jwt.Secure)

		auth := cookies[1]
		assert.Equal(t, IsAuthenticatedCookie, auth.Name)
		assert.Equal(t, "true", auth.Value)
		assert.False(t, auth.HttpOnly)
	})

	t.Run("ClearAuthCookies", func(t *testing.T) {
		w := httptest.NewRecorder()
		ClearAuthCookies(w, false)

		cookies := w.Result().Cookies()
		assert.Len(t, cookies, 2)
		assert.Equal(t, "", cookies[0].Value)
		assert.True(t, cookies[0].Expires.Before(time.Now()))
	})
}

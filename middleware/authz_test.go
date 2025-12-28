package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizeRole(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Authorized User - Admin", func(t *testing.T) {
		user := &middleware.AuthenticatedUser{
			ID:   "1",
			Role: "admin",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler := middleware.AuthorizeRole("admin", "editor")(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Authorized User - Editor", func(t *testing.T) {
		user := &middleware.AuthenticatedUser{
			ID:   "2",
			Role: "editor",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler := middleware.AuthorizeRole("admin", "editor")(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Unauthorized User - Wrong Role", func(t *testing.T) {
		user := &middleware.AuthenticatedUser{
			ID:   "3",
			Role: "user",
		}
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), middleware.UserContextKey, user)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler := middleware.AuthorizeRole("admin", "editor")(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Contains(t, rr.Body.String(), "forbidden: insufficient permissions")
	})

	t.Run("Unauthorized User - No User in Context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler := middleware.AuthorizeRole("admin")(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "unauthorized: user not found in context")
	})
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("User exists", func(t *testing.T) {
		expectedUser := &middleware.AuthenticatedUser{ID: "123", Role: "user"}
		ctx := context.WithValue(context.Background(), middleware.UserContextKey, expectedUser)
		
		user, ok := middleware.GetUserFromContext(ctx)
		
		assert.True(t, ok)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("User does not exist", func(t *testing.T) {
		user, ok := middleware.GetUserFromContext(context.Background())
		
		assert.False(t, ok)
		assert.Nil(t, user)
	})
}

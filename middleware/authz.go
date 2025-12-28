package middleware

import (
	"context"
	"net/http"

	customErrors "github.com/shashtag-ventures/go-common/errors"
	"github.com/shashtag-ventures/go-common/jsonResponse"
)

// AuthorizeRole checks if the authenticated user has one of the required roles.
func AuthorizeRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			user, ok := r.Context().Value(UserContextKey).(*AuthenticatedUser)
			if !ok || user == nil {
				jsonResponse.SendErrorResponse(w, customErrors.New("unauthorized: user not found in context", nil), http.StatusUnauthorized)
				return
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, requiredRole := range requiredRoles {
				if user.Role == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				jsonResponse.SendErrorResponse(w, customErrors.New("forbidden: insufficient permissions", nil), http.StatusForbidden)
				return
			}

			// If authorized, pass the request to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// AuthenticatedUser represents the authenticated user's information.
type AuthenticatedUser struct {
	ID    string
	Email string
	Role  string
}

// GetUserFromContext extracts the AuthenticatedUser from the request context.
func GetUserFromContext(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(UserContextKey).(*AuthenticatedUser)
	return user, ok
}

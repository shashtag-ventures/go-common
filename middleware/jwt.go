package middleware

import (
	"context"
	"fmt"
	"net/http"

	customErrors "github.com/shashtag-ventures/go-common/errors"
	"github.com/shashtag-ventures/go-common/jsonResponse"
	"github.com/shashtag-ventures/go-common/jwt"
)

// JWTAuthMiddleware creates a middleware that authenticates requests using a JWT from a cookie.
func JWTAuthMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt_token")
			if err != nil {
				jsonResponse.SendErrorResponse(w, customErrors.New("missing jwt cookie", err), http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value

			claims, err := jwt.ParseToken(tokenString, jwtSecret)
			if err != nil {
				jsonResponse.SendErrorResponse(w, fmt.Errorf("invalid or expired token: %w", err), http.StatusUnauthorized)
				return
			}

			authenticatedUser := &AuthenticatedUser{
				ID:    claims.UserID,
				Email: "", // Email is not in claims
				Role:  claims.Role,
			}

			// NEW: Capture UserID in the mutable log state for the outer logger
			if state, ok := r.Context().Value(LogStateKey).(*LogState); ok {
				state.SetUser(authenticatedUser.ID)
			}

			ctx := context.WithValue(r.Context(), UserContextKey, authenticatedUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

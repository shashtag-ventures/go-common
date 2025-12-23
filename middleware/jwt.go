package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	customErrors "github.com/shashtag-ventures/go-common/errors"
	"github.com/shashtag-ventures/go-common/jsonResponse"
	"github.com/shashtag-ventures/go-common/jwt"
)

// JWTAuthMiddleware creates a middleware that authenticates requests using a JWT from a cookie.
func JWTAuthMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("JWTAuthMiddleware: Request received", "method", r.Method, "url", r.URL.String())
			cookie, err := r.Cookie("jwt_token")
			if err != nil {
				slog.Warn("JWTAuthMiddleware: Missing JWT cookie", "error", err)
				jsonResponse.SendErrorResponse(w, customErrors.New("missing jwt cookie", err), http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value

			claims, err := jwt.ParseToken(tokenString, jwtSecret)
			if err != nil {
				slog.Error("JWTAuthMiddleware: Invalid or expired token", "error", err)
				jsonResponse.SendErrorResponse(w, fmt.Errorf("invalid or expired token: %w", err), http.StatusUnauthorized)
				return
			}

			// The AuthenticatedUser struct is defined in authz.go
			authenticatedUser := &AuthenticatedUser{
				ID:    claims.UserID,
				Email: "", // Email is not in claims
				Role:  claims.Role, // Role is now a simple string
			}

			slog.Info("JWTAuthMiddleware: User authenticated from token", "userID", authenticatedUser.ID, "role", authenticatedUser.Role)

			ctx := context.WithValue(r.Context(), UserContextKey, authenticatedUser)
			slog.Info("JWTAuthMiddleware: User set in context")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

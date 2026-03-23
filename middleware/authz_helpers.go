package middleware

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/errors"
)

// GetAuthenticatedUser retrieves the AuthenticatedUser from the context or returns an error.
func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, error) {
	user, ok := GetUserFromContext(ctx)
	if !ok || user == nil {
		return nil, errors.New("unauthorized: user not found in context", errors.ErrUnauthorized)
	}
	return user, nil
}

// GetAuthenticatedUserID retrieves and parses the user's UUID from the context.
func GetAuthenticatedUserID(ctx context.Context) (uuid.UUID, error) {
	user, err := GetAuthenticatedUser(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(user.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user id format: %w", err)
	}

	return id, nil
}

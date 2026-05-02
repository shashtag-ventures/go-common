package authz

import (
	"context"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/errors"
)

// MembershipStore defines the interface for checking resource membership.
// Any service that needs to verify user access to a resource (team, project, etc.)
// should implement this interface or provide a store that does.
type MembershipStore interface {
	IsMember(ctx context.Context, resourceID, userID uuid.UUID) (bool, error)
}

// AdminMembershipStore defines the interface for checking admin-level access to a resource.
type AdminMembershipStore interface {
	IsAdmin(ctx context.Context, resourceID, userID uuid.UUID) (bool, error)
}

// CheckMembership verifies if a user has access to a resource via a MembershipStore.
// It returns a wrapped errors.ErrForbidden if the user is not a member.
func CheckMembership(ctx context.Context, store MembershipStore, resourceID, userID uuid.UUID) error {
	isMember, err := store.IsMember(ctx, resourceID, userID)
	if err != nil {
		return errors.New("failed to verify membership", err)
	}
	if !isMember {
		return errors.New("user is not a member of this resource", errors.ErrForbidden)
	}
	return nil
}

// CheckAdminMembership verifies if a user has admin access to a resource.
// It returns a wrapped errors.ErrForbidden if the user is not an admin.
func CheckAdminMembership(ctx context.Context, store AdminMembershipStore, resourceID, userID uuid.UUID) error {
	isAdmin, err := store.IsAdmin(ctx, resourceID, userID)
	if err != nil {
		return errors.New("failed to verify admin membership", err)
	}
	if !isAdmin {
		return errors.New("user is not an admin of this resource", errors.ErrForbidden)
	}
	return nil
}

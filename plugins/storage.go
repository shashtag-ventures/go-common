package plugins

import (
	"context"

	"github.com/google/uuid"
)

// PluginStorage defines the interface for storing external connections.
type PluginStorage interface {
	SaveConnection(ctx context.Context, conn *ExternalConnection) error
	GetConnection(ctx context.Context, userID uuid.UUID, provider string) (*ExternalConnection, error)
	GetConnectionByProviderID(ctx context.Context, provider string, providerUserID string) (*ExternalConnection, error)
	ListConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error)
	UpdateInstallationID(ctx context.Context, userID uuid.UUID, provider string, installationID string) error
	DeleteConnection(ctx context.Context, userID uuid.UUID, provider string) error
}

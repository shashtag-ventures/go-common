package integrations

import (
	"context"
)

// IntegrationStorage defines the interface for storing external connections.
type IntegrationStorage interface {
	SaveConnection(ctx context.Context, conn *ExternalConnection) error
	GetConnection(ctx context.Context, userID uint, provider string) (*ExternalConnection, error)
	ListConnections(ctx context.Context, userID uint) ([]*ExternalConnection, error)
}

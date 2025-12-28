package integrations

import (
	"context"

	"github.com/google/uuid"
)

// IntegrationStorage defines the interface for storing external connections.
type IntegrationStorage interface {
	SaveConnection(ctx context.Context, conn *ExternalConnection) error
	GetConnection(ctx context.Context, userID uuid.UUID, provider string) (*ExternalConnection, error)
	ListConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error)
}

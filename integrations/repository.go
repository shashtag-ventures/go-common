package integrations

import (
	"context"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/gormutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type integrationRepository struct {
	repo *gormutil.Repository[ExternalConnection]
}

func NewIntegrationRepository(db *gorm.DB) IntegrationStorage {
	return &integrationRepository{
		repo: gormutil.NewRepository[ExternalConnection](db),
	}
}

func (r *integrationRepository) SaveConnection(ctx context.Context, conn *ExternalConnection) error {
	// Upsert: Create or Update on conflict of (user_id, provider)
	return r.repo.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "provider"}},
		DoUpdates: clause.AssignmentColumns([]string{"access_token", "refresh_token", "expires_at", "username", "avatar_url", "provider_user_id", "installation_id"}),
	}).Create(conn).Error
}

func (r *integrationRepository) GetConnection(ctx context.Context, userID uuid.UUID, provider string) (*ExternalConnection, error) {
	return r.repo.FindOneBy(ctx, "user_id = ? AND provider = ?", userID, provider)
}

func (r *integrationRepository) GetConnectionByProviderID(ctx context.Context, provider string, providerUserID string) (*ExternalConnection, error) {
	return r.repo.FindOneBy(ctx, "provider = ? AND provider_user_id = ?", provider, providerUserID)
}

func (r *integrationRepository) ListConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error) {
	return r.repo.Find(ctx, "user_id = ?", userID)
}

func (r *integrationRepository) UpdateInstallationID(ctx context.Context, userID uuid.UUID, provider string, installationID string) error {
	conn := &ExternalConnection{
		UserID:         userID,
		Provider:       provider,
		InstallationID: installationID,
	}
	// Upsert: Create or update only the installation_id field
	return r.repo.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "provider"}},
		DoUpdates: clause.AssignmentColumns([]string{"installation_id"}),
	}).Create(conn).Error
}

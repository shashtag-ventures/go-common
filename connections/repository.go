package connections

import (
	"context"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/gormutil"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type connectionRepository struct {
	repo *gormutil.Repository[ExternalConnection]
}

func NewConnectionRepository(db *gorm.DB) ConnectionStorage {
	return &connectionRepository{
		repo: gormutil.NewRepository[ExternalConnection](db),
	}
}

func (r *connectionRepository) SaveConnection(ctx context.Context, conn *ExternalConnection) error {
	// Upsert: Create or Update on conflict of (user_id, provider).
	// Note: installation_id is intentionally excluded from DoUpdates — it is managed
	// exclusively by UpdateInstallationID (GitHub App setup flow). Including it here
	// would cause OAuth callbacks (which don't carry installation_id) to wipe
	// any previously saved value with an empty string.
	return r.repo.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "provider"}},
		DoUpdates: clause.AssignmentColumns([]string{"access_token", "refresh_token", "expires_at", "username", "avatar_url", "provider_user_id"}),
	}).Create(conn).Error
}

func (r *connectionRepository) GetConnection(ctx context.Context, userID uuid.UUID, provider string) (*ExternalConnection, error) {
	return r.repo.FindOneBy(ctx, "user_id = ? AND provider = ?", userID, provider)
}

func (r *connectionRepository) GetConnectionByProviderID(ctx context.Context, provider string, providerUserID string) (*ExternalConnection, error) {
	return r.repo.FindOneBy(ctx, "provider = ? AND provider_user_id = ?", provider, providerUserID)
}

func (r *connectionRepository) ListConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error) {
	return r.repo.Find(ctx, "user_id = ?", userID)
}

func (r *connectionRepository) UpdateInstallationID(ctx context.Context, userID uuid.UUID, provider string, installationID string) error {
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

func (r *connectionRepository) DeleteConnection(ctx context.Context, userID uuid.UUID, provider string) error {
	return r.repo.DB(ctx).Where("user_id = ? AND provider = ?", userID, provider).Delete(&ExternalConnection{}).Error
}

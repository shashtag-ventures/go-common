package integrations

import (
	"time"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/gormutil"
)

// ExternalConnection represents a link between a local user and an external OAuth provider.
type ExternalConnection struct {
	gormutil.BaseModel
	UserID           uuid.UUID `json:"user_id" gorm:"type:uuid;uniqueIndex:idx_user_provider"`
	Provider         string    `json:"provider" gorm:"uniqueIndex:idx_user_provider"` // e.g., "github", "google"
	ProviderUserID   string    `json:"provider_user_id"`                             // The ID assigned by the provider
	AccessToken      string    `json:"-"`                                            // Should be encrypted
	RefreshToken     string    `json:"-"`                                            // Should be encrypted
	ExpiresAt        time.Time `json:"expires_at"`
	Username         string    `json:"username"`   // The handle on the provider
	AvatarURL        string    `json:"avatar_url"`
}

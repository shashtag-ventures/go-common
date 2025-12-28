package integrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/crypto"
	"github.com/shashtag-ventures/go-common/integrations/clients"
	"github.com/shashtag-ventures/go-common/integrations/types"
	"github.com/shashtag-ventures/go-common/middleware"
)

// IntegrationService defines operations for managing external integrations.
type IntegrationService interface {
	SaveConnection(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, username string, avatarURL string, accessToken string, refreshToken string) error
	GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error)
	ListUserRepositories(ctx context.Context, userID uuid.UUID, provider string) ([]types.Repository, error)
	ListUserNamespaces(ctx context.Context, userID uuid.UUID, provider string) ([]types.Namespace, error)
}

type integrationService struct {
	db                 IntegrationStorage
	tokenEncryptionKey string
	githubClient       *clients.GitHubClient
}

// NewIntegrationService creates a new instance of IntegrationService.
func NewIntegrationService(db IntegrationStorage, tokenEncryptionKey string) IntegrationService {
	return &integrationService{
		db:                 db,
		tokenEncryptionKey: tokenEncryptionKey,
		githubClient:       clients.NewGitHubClient(),
	}
}

func (s *integrationService) SaveConnection(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, username string, avatarURL string, accessToken string, refreshToken string) error {
	logger := middleware.GetLoggerFromContext(ctx)

	// Encrypt tokens before saving
	encryptedAccess, err := crypto.Encrypt(accessToken, s.tokenEncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	var encryptedRefresh string
	if refreshToken != "" {
		encryptedRefresh, err = crypto.Encrypt(refreshToken, s.tokenEncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt refresh token: %w", err)
		}
	}

	conn := &ExternalConnection{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		AccessToken:    encryptedAccess,
		RefreshToken:   encryptedRefresh,
		Username:       username,
		AvatarURL:      avatarURL,
	}

	if err := s.db.SaveConnection(ctx, conn); err != nil {
		logger.Error("Failed to save integration connection", "userID", userID, "provider", provider, "error", err)
		return err
	}

	return nil
}

func (s *integrationService) GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error) {
	return s.db.ListConnections(ctx, userID)
}

func (s *integrationService) ListUserRepositories(ctx context.Context, userID uuid.UUID, provider string) ([]types.Repository, error) {
	conn, err := s.db.GetConnection(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	accessToken, err := crypto.Decrypt(conn.AccessToken, s.tokenEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	switch provider {
	case "github":
		return s.githubClient.ListRepositories(ctx, accessToken)
	default:
		return nil, fmt.Errorf("provider %s not supported for repo listing", provider)
	}
}

func (s *integrationService) ListUserNamespaces(ctx context.Context, userID uuid.UUID, provider string) ([]types.Namespace, error) {
	conn, err := s.db.GetConnection(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	accessToken, err := crypto.Decrypt(conn.AccessToken, s.tokenEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt access token: %w", err)
	}

	switch provider {
	case "github":
		return s.githubClient.ListNamespaces(ctx, accessToken)
	default:
		return nil, fmt.Errorf("provider %s not supported for namespace listing", provider)
	}
}

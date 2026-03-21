package integrations

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/crypto"
	"github.com/shashtag-ventures/go-common/integrations/types"
	"github.com/shashtag-ventures/go-common/middleware"
	"time"
)

type IntegrationService interface {
	SaveConnection(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, username string, avatarURL string, accessToken string, refreshToken string, expiresAt time.Time) error
	GetConnectionByProviderID(ctx context.Context, provider string, providerUserID string) (*ExternalConnection, error)
	GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error)
	ListUserRepositories(ctx context.Context, userID uuid.UUID, provider string) ([]types.Repository, error)
	ListUserNamespaces(ctx context.Context, userID uuid.UUID, provider string) ([]types.Namespace, error)
}

type integrationService struct {
	db                 IntegrationStorage
	tokenEncryptionKey string
	clients            map[string]types.IntegrationClient
}

// NewIntegrationService creates a new instance of IntegrationService.
func NewIntegrationService(db IntegrationStorage, tokenEncryptionKey string, clients map[string]types.IntegrationClient) IntegrationService {
	return &integrationService{
		db:                 db,
		tokenEncryptionKey: tokenEncryptionKey,
		clients:            clients,
	}
}

func (s *integrationService) SaveConnection(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, username string, avatarURL string, accessToken string, refreshToken string, expiresAt time.Time) error {
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
		ExpiresAt:      expiresAt,
		Username:       username,
		AvatarURL:      avatarURL,
	}

	if err := s.db.SaveConnection(ctx, conn); err != nil {
		logger.Error("Failed to save integration connection", "userID", userID, "provider", provider, "error", err)
		return err
	}

	return nil
}

func (s *integrationService) GetConnectionByProviderID(ctx context.Context, provider string, providerUserID string) (*ExternalConnection, error) {
	return s.db.GetConnectionByProviderID(ctx, provider, providerUserID)
}

func (s *integrationService) GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error) {
	return s.db.ListConnections(ctx, userID)
}

func (s *integrationService) ensureValidToken(ctx context.Context, conn *ExternalConnection, client types.IntegrationClient) (string, error) {
	accessToken, err := crypto.Decrypt(conn.AccessToken, s.tokenEncryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt access token: %w", err)
	}

	if !conn.ExpiresAt.IsZero() && time.Now().Add(1*time.Minute).After(conn.ExpiresAt) {
		decryptedRefresh, err := crypto.Decrypt(conn.RefreshToken, s.tokenEncryptionKey)
		if err == nil && decryptedRefresh != "" {
			refreshResp, err := client.RefreshToken(ctx, decryptedRefresh)
			if err == nil && refreshResp != nil {
				accessToken = refreshResp.AccessToken
				
				conn.AccessToken, _ = crypto.Encrypt(refreshResp.AccessToken, s.tokenEncryptionKey)
				if refreshResp.RefreshToken != "" {
					conn.RefreshToken, _ = crypto.Encrypt(refreshResp.RefreshToken, s.tokenEncryptionKey)
				}
				if !refreshResp.ExpiresAt.IsZero() {
					conn.ExpiresAt = refreshResp.ExpiresAt
				}
				
				if err := s.db.SaveConnection(ctx, conn); err != nil {
					middleware.GetLoggerFromContext(ctx).Error("Failed to save refreshed token", "error", err)
				}
			} else {
				middleware.GetLoggerFromContext(ctx).Error("Failed to refresh token", "error", err)
			}
		}
	}

	return accessToken, nil
}

func (s *integrationService) ListUserRepositories(ctx context.Context, userID uuid.UUID, provider string) ([]types.Repository, error) {
	conn, err := s.db.GetConnection(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	client, ok := s.clients[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not supported for repository listing", provider)
	}

	accessToken, err := s.ensureValidToken(ctx, conn, client)
	if err != nil {
		return nil, err
	}

	return client.ListRepositories(ctx, accessToken)
}

func (s *integrationService) ListUserNamespaces(ctx context.Context, userID uuid.UUID, provider string) ([]types.Namespace, error) {
	conn, err := s.db.GetConnection(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	client, ok := s.clients[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not supported for namespace listing", provider)
	}

	accessToken, err := s.ensureValidToken(ctx, conn, client)
	if err != nil {
		return nil, err
	}

	return client.ListNamespaces(ctx, accessToken)
}

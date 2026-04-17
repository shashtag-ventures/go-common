package plugins

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/crypto"
	"github.com/shashtag-ventures/go-common/plugins/types"
	"github.com/shashtag-ventures/go-common/middleware"
	"time"
)

type PluginService interface {
	SaveConnection(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, username string, avatarURL string, accessToken string, refreshToken string, expiresAt time.Time, installationID string) error
	SaveInstallation(ctx context.Context, userID uuid.UUID, provider string, installationID string) error
	GetConnection(ctx context.Context, userID uuid.UUID, provider string) (*ExternalConnection, error)
	GetConnectionByProviderID(ctx context.Context, provider string, providerUserID string) (*ExternalConnection, error)
	GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error)
	ListUserRepositories(ctx context.Context, userID uuid.UUID, provider string) ([]types.Repository, error)
	ListUserRepositoriesPaginated(ctx context.Context, userID uuid.UUID, provider string, search string, namespace string, page int, limit int) ([]types.Repository, error)
	ListUserNamespaces(ctx context.Context, userID uuid.UUID, provider string) ([]types.Namespace, error)
	ListRepositoryContents(ctx context.Context, userID uuid.UUID, provider string, repoFullName string, path string) ([]types.ContentItem, error)
	DeleteConnection(ctx context.Context, userID uuid.UUID, provider string) error
}

type pluginService struct {
	db                 PluginStorage
	tokenEncryptionKey string
	clients            map[string]types.PluginClient
}

// NewPluginService creates a new instance of PluginService.
func NewPluginService(db PluginStorage, tokenEncryptionKey string, clients map[string]types.PluginClient) PluginService {
	return &pluginService{
		db:                 db,
		tokenEncryptionKey: tokenEncryptionKey,
		clients:            clients,
	}
}

func (s *pluginService) SaveConnection(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, username string, avatarURL string, accessToken string, refreshToken string, expiresAt time.Time, installationID string) error {
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
		InstallationID: installationID,
	}

	if err := s.db.SaveConnection(ctx, conn); err != nil {
		logger.Error("Failed to save plugin connection", "userID", userID, "provider", provider, "error", err)
		return err
	}

	return nil
}

func (s *pluginService) GetConnection(ctx context.Context, userID uuid.UUID, provider string) (*ExternalConnection, error) {
	return s.db.GetConnection(ctx, userID, provider)
}

func (s *pluginService) GetConnectionByProviderID(ctx context.Context, provider string, providerUserID string) (*ExternalConnection, error) {
	return s.db.GetConnectionByProviderID(ctx, provider, providerUserID)
}

func (s *pluginService) GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error) {
	return s.db.ListConnections(ctx, userID)
}

func (s *pluginService) SaveInstallation(ctx context.Context, userID uuid.UUID, provider string, installationID string) error {
	logger := middleware.GetLoggerFromContext(ctx)
	if err := s.db.UpdateInstallationID(ctx, userID, provider, installationID); err != nil {
		logger.Error("Failed to save installation", "userID", userID, "provider", provider, "installationID", installationID, "error", err)
		return err
	}
	return nil
}

func (s *pluginService) ensureValidToken(ctx context.Context, conn *ExternalConnection, client types.PluginClient) (string, error) {
	if conn.AccessToken == "" {
		// If AccessToken is empty, it means we don't have an OAuth token.
		// We return an empty string (no error) so that the caller can proceed
		// and use the InstallationID if available.
		return "", nil
	}

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

func (s *pluginService) ListUserRepositories(ctx context.Context, userID uuid.UUID, provider string) ([]types.Repository, error) {
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

	return client.ListRepositories(ctx, accessToken, conn.InstallationID)
}

func (s *pluginService) ListUserRepositoriesPaginated(ctx context.Context, userID uuid.UUID, provider string, search string, namespace string, page int, limit int) ([]types.Repository, error) {
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

	if search != "" {
		if namespace == "" || namespace == "all" {
			namespace = conn.Username
		}
		return client.SearchRepositories(ctx, accessToken, search, namespace, page, limit, conn.InstallationID)
	}

	return client.ListRepositoriesPaginated(ctx, accessToken, conn.InstallationID, page, limit)
}

func (s *pluginService) ListUserNamespaces(ctx context.Context, userID uuid.UUID, provider string) ([]types.Namespace, error) {
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

	return client.ListNamespaces(ctx, accessToken, conn.InstallationID)
}

func (s *pluginService) ListRepositoryContents(ctx context.Context, userID uuid.UUID, provider string, repoFullName string, path string) ([]types.ContentItem, error) {
	conn, err := s.db.GetConnection(ctx, userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	client, ok := s.clients[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not supported for content listing", provider)
	}

	accessToken, err := s.ensureValidToken(ctx, conn, client)
	if err != nil {
		return nil, err
	}

	return client.ListContents(ctx, accessToken, repoFullName, path, conn.InstallationID)
}

func (s *pluginService) DeleteConnection(ctx context.Context, userID uuid.UUID, provider string) error {
	logger := middleware.GetLoggerFromContext(ctx)
	if err := s.db.DeleteConnection(ctx, userID, provider); err != nil {
		logger.Error("Failed to delete plugin connection", "userID", userID, "provider", provider, "error", err)
		return err
	}
	logger.Info("Plugin connection deleted", "userID", userID, "provider", provider)
	return nil
}

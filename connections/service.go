package connections

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/crypto"
	"github.com/shashtag-ventures/go-common/connections/types"
	"time"
)

// LoggerFunc extracts a logger from context. This allows callers to inject
// their own logger strategy (e.g. middleware.GetLoggerFromContext) without
// coupling this package to any specific middleware implementation.
type LoggerFunc func(context.Context) *slog.Logger

// SaveConnectionParams holds all data needed to save or update an external connection.
type SaveConnectionParams struct {
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
	Username       string
	AvatarURL      string
	AccessToken    string
	RefreshToken   string
	ExpiresAt      time.Time
	InstallationID string
}

type ConnectionService interface {
	SaveConnection(ctx context.Context, params SaveConnectionParams) error
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

type connectionService struct {
	db                 ConnectionStorage
	tokenEncryptionKey string
	clients            map[string]types.ProviderClient
	getLogger          LoggerFunc
}

// NewConnectionService creates a new instance of ConnectionService.
// logFn is used to extract a logger from context for structured logging.
func NewConnectionService(db ConnectionStorage, tokenEncryptionKey string, clients map[string]types.ProviderClient, logFn LoggerFunc) ConnectionService {
	if logFn == nil {
		logFn = func(_ context.Context) *slog.Logger { return slog.Default() }
	}
	return &connectionService{
		db:                 db,
		tokenEncryptionKey: tokenEncryptionKey,
		clients:            clients,
		getLogger:          logFn,
	}
}

func (s *connectionService) SaveConnection(ctx context.Context, params SaveConnectionParams) error {
	logger := s.getLogger(ctx)

	// Encrypt tokens before saving
	encryptedAccess, err := crypto.Encrypt(params.AccessToken, s.tokenEncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt access token: %w", err)
	}

	var encryptedRefresh string
	if params.RefreshToken != "" {
		encryptedRefresh, err = crypto.Encrypt(params.RefreshToken, s.tokenEncryptionKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt refresh token: %w", err)
		}
	}

	conn := &ExternalConnection{
		UserID:         params.UserID,
		Provider:       params.Provider,
		ProviderUserID: params.ProviderUserID,
		AccessToken:    encryptedAccess,
		RefreshToken:   encryptedRefresh,
		ExpiresAt:      params.ExpiresAt,
		Username:       params.Username,
		AvatarURL:      params.AvatarURL,
		InstallationID: params.InstallationID,
	}

	if err := s.db.SaveConnection(ctx, conn); err != nil {
		logger.Error("Failed to save connection", "userID", params.UserID, "provider", params.Provider, "error", err)
		return err
	}

	return nil
}

func (s *connectionService) GetConnection(ctx context.Context, userID uuid.UUID, provider string) (*ExternalConnection, error) {
	return s.db.GetConnection(ctx, userID, provider)
}

func (s *connectionService) GetConnectionByProviderID(ctx context.Context, provider string, providerUserID string) (*ExternalConnection, error) {
	return s.db.GetConnectionByProviderID(ctx, provider, providerUserID)
}

func (s *connectionService) GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ExternalConnection, error) {
	return s.db.ListConnections(ctx, userID)
}

func (s *connectionService) SaveInstallation(ctx context.Context, userID uuid.UUID, provider string, installationID string) error {
	logger := s.getLogger(ctx)
	if err := s.db.UpdateInstallationID(ctx, userID, provider, installationID); err != nil {
		logger.Error("Failed to save installation", "userID", userID, "provider", provider, "installationID", installationID, "error", err)
		return err
	}
	return nil
}

func (s *connectionService) ensureValidToken(ctx context.Context, conn *ExternalConnection, client types.ProviderClient) (string, error) {
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

				encryptedAccess, err := crypto.Encrypt(refreshResp.AccessToken, s.tokenEncryptionKey)
				if err != nil {
					s.getLogger(ctx).Error("Failed to encrypt refreshed access token", "error", err)
					return accessToken, nil
				}
				conn.AccessToken = encryptedAccess

				if refreshResp.RefreshToken != "" {
					encryptedRefresh, err := crypto.Encrypt(refreshResp.RefreshToken, s.tokenEncryptionKey)
					if err != nil {
						s.getLogger(ctx).Error("Failed to encrypt refreshed refresh token", "error", err)
						return accessToken, nil
					}
					conn.RefreshToken = encryptedRefresh
				}
				if !refreshResp.ExpiresAt.IsZero() {
					conn.ExpiresAt = refreshResp.ExpiresAt
				}

				if err := s.db.SaveConnection(ctx, conn); err != nil {
					s.getLogger(ctx).Error("Failed to save refreshed token", "error", err)
				}
			} else {
				s.getLogger(ctx).Error("Failed to refresh token", "error", err)
			}
		}
	}

	return accessToken, nil
}

func (s *connectionService) ListUserRepositories(ctx context.Context, userID uuid.UUID, provider string) ([]types.Repository, error) {
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

	// Use OAuth token path (empty installationID) so the user can see repos
	// across ALL orgs where the GitHub App is installed, not just one.
	return client.ListRepositories(ctx, accessToken, "")
}

func (s *connectionService) ListUserRepositoriesPaginated(ctx context.Context, userID uuid.UUID, provider string, search string, namespace string, page int, limit int) ([]types.Repository, error) {
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
		return client.SearchRepositories(ctx, accessToken, search, namespace, page, limit, "")
	}

	return client.ListRepositoriesPaginated(ctx, accessToken, "", page, limit)
}

func (s *connectionService) ListUserNamespaces(ctx context.Context, userID uuid.UUID, provider string) ([]types.Namespace, error) {
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

	// Use OAuth token path (empty installationID) so the user sees ALL orgs
	// where the GitHub App is installed, not just the single stored installation.
	return client.ListNamespaces(ctx, accessToken, "")
}

func (s *connectionService) ListRepositoryContents(ctx context.Context, userID uuid.UUID, provider string, repoFullName string, path string) ([]types.ContentItem, error) {
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

func (s *connectionService) DeleteConnection(ctx context.Context, userID uuid.UUID, provider string) error {
	logger := s.getLogger(ctx)
	if err := s.db.DeleteConnection(ctx, userID, provider); err != nil {
		logger.Error("Failed to delete connection", "userID", userID, "provider", provider, "error", err)
		return err
	}
	logger.Info("Connection deleted", "userID", userID, "provider", provider)
	return nil
}

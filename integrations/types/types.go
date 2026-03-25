package types

import (
	"context"
	"time"
)

type TokenRefreshResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"-"`
}

type IntegrationClient interface {
	ListRepositories(ctx context.Context, token string, installationID string) ([]Repository, error)
	ListRepositoriesPaginated(ctx context.Context, token string, installationID string, page int, limit int) ([]Repository, error)
	SearchRepositories(ctx context.Context, token string, query string, namespace string, page int, limit int, installationID string) ([]Repository, error)
	ListNamespaces(ctx context.Context, token string, installationID string) ([]Namespace, error)
	ListContents(ctx context.Context, token string, repoFullName string, path string, installationID string) ([]ContentItem, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenRefreshResponse, error)
}

type ContentItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"` // "file" or "dir"
	Size int64  `json:"size"`
}

type Repository struct {
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	URL       string    `json:"url"`
	Private   bool      `json:"private"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Namespace struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"` // "User" or "Organization"
}

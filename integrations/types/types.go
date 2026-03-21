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
	ListRepositories(ctx context.Context, token string) ([]Repository, error)
	ListRepositoriesPaginated(ctx context.Context, token string, page int, limit int) ([]Repository, error)
	SearchRepositories(ctx context.Context, token string, query string, namespace string, page int, limit int) ([]Repository, error)
	ListNamespaces(ctx context.Context, token string) ([]Namespace, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenRefreshResponse, error)
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

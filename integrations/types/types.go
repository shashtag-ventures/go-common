package types

import (
	"context"
	"time"
)

type IntegrationClient interface {
	ListRepositories(ctx context.Context, token string) ([]Repository, error)
	ListNamespaces(ctx context.Context, token string) ([]Namespace, error)
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

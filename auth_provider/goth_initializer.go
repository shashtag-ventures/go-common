package auth_provider

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

// GothConfig holds the configuration for the Goth initializer.
type GothConfig struct {
	SessionSecret      string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string
	GitHubClientID     string
	GitHubClientSecret string
	GitHubCallbackURL  string
}

// OAuthProviderInitializer defines the interface for initializing OAuth providers.
type OAuthProviderInitializer interface {
	Init() error
}

// GothInitializer implements OAuthProviderInitializer for Goth (Go OAuth).
type GothInitializer struct {
	cfg GothConfig
}

// NewGothInitializer creates a new instance of GothInitializer.
func NewGothInitializer(cfg GothConfig) *GothInitializer {
	return &GothInitializer{
		cfg: cfg,
	}
}

// Init initializes Goth with configured OAuth providers and session store.
func (g *GothInitializer) Init() error {
	// Use Gorilla Sessions for storing Goth sessions.
	gothic.Store = sessions.NewCookieStore([]byte(g.cfg.SessionSecret))

	var providers []goth.Provider

	// Register Google OAuth2 provider if credentials are provided.
	if g.cfg.GoogleClientID != "" && g.cfg.GoogleClientSecret != "" {
		providers = append(providers, google.New(g.cfg.GoogleClientID, g.cfg.GoogleClientSecret, g.cfg.GoogleCallbackURL, "email", "profile"))
	}

	// Register GitHub OAuth2 provider if credentials are provided.
	if g.cfg.GitHubClientID != "" && g.cfg.GitHubClientSecret != "" {
		providers = append(providers, github.New(g.cfg.GitHubClientID, g.cfg.GitHubClientSecret, g.cfg.GitHubCallbackURL))
	}

	goth.UseProviders(providers...)
	return nil
}

package auth_provider

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// GothConfig holds the configuration for the Goth initializer.
type GothConfig struct {
	SessionSecret      string
	GoogleClientID     string
	GoogleClientSecret string
	CallbackURL        string
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

	// Register Google OAuth2 provider.
	goth.UseProviders(
		google.New(g.cfg.GoogleClientID, g.cfg.GoogleClientSecret, g.cfg.CallbackURL, "email", "profile"),
	)
	return nil
}

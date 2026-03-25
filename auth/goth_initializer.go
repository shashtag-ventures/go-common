package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
	"net/http"
)

// GothProvider defines the interface for Goth authentication methods to allow for easier testing.
type GothProvider interface {
	CompleteUserAuth(w http.ResponseWriter, r *http.Request) (goth.User, error)
	BeginAuthHandler(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request) error
	GetSession(r *http.Request, name string) (*sessions.Session, error)
}

// RealGothProvider is the production implementation of GothProvider that calls the goth/gothic functions.
type RealGothProvider struct{}

func (p *RealGothProvider) CompleteUserAuth(w http.ResponseWriter, r *http.Request) (goth.User, error) {
	return gothic.CompleteUserAuth(w, r)
}

func (p *RealGothProvider) BeginAuthHandler(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func (p *RealGothProvider) Logout(w http.ResponseWriter, r *http.Request) error {
	return gothic.Logout(w, r)
}

func (p *RealGothProvider) GetSession(r *http.Request, name string) (*sessions.Session, error) {
	return gothic.Store.Get(r, name)
}

// GothConfig holds the configuration for the Goth initializer.
type GothConfig struct {
	SessionSecret         string
	GoogleClientID        string
	GoogleClientSecret    string
	GoogleCallbackURL     string
	GitHubClientID        string
	GitHubClientSecret    string
	GitHubCallbackURL     string
	GitLabClientID        string
	GitLabClientSecret    string
	GitLabCallbackURL     string
	BitbucketClientID     string
	BitbucketClientSecret string
	BitbucketCallbackURL  string
	MicrosoftClientID     string
	MicrosoftClientSecret string
	MicrosoftCallbackURL  string
	CookieDomain          string
	Secure                bool
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
func NewGothInitializer(cfg GothConfig) *gothInitializer {
	return &gothInitializer{
		cfg: cfg,
	}
}

// gothInitializer is the internal implementation of OAuthProviderInitializer.
type gothInitializer struct {
	cfg GothConfig
}

// Init initializes Goth with configured OAuth providers and session store.
func (g *gothInitializer) Init() error {
	// Use Gorilla Sessions for storing Goth sessions.
	store := sessions.NewCookieStore([]byte(g.cfg.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   g.cfg.CookieDomain,
		MaxAge:   86400 * 30, // 30 days
		HttpOnly: true,
		Secure:   g.cfg.Secure,
		SameSite: 2, // http.SameSiteLaxMode is 2
	}
	gothic.Store = store

	var providers []goth.Provider

	// Register Google OAuth2 provider if credentials are provided.
	if g.cfg.GoogleClientID != "" && g.cfg.GoogleClientSecret != "" {
		providers = append(providers, google.New(g.cfg.GoogleClientID, g.cfg.GoogleClientSecret, g.cfg.GoogleCallbackURL, "email", "profile"))
	}

	// Register GitHub OAuth2 provider if credentials are provided.
	if g.cfg.GitHubClientID != "" && g.cfg.GitHubClientSecret != "" {
		providers = append(providers, github.New(g.cfg.GitHubClientID, g.cfg.GitHubClientSecret, g.cfg.GitHubCallbackURL, "user:email"))
	}

	// Register GitLab OAuth2 provider if credentials are provided.
	if g.cfg.GitLabClientID != "" && g.cfg.GitLabClientSecret != "" {
		providers = append(providers, gitlab.New(g.cfg.GitLabClientID, g.cfg.GitLabClientSecret, g.cfg.GitLabCallbackURL))
	}

	// Register Bitbucket OAuth2 provider if credentials are provided.
	if g.cfg.BitbucketClientID != "" && g.cfg.BitbucketClientSecret != "" {
		providers = append(providers, bitbucket.New(g.cfg.BitbucketClientID, g.cfg.BitbucketClientSecret, g.cfg.BitbucketCallbackURL))
	}

	// Register Microsoft OAuth2 provider if credentials are provided.
	if g.cfg.MicrosoftClientID != "" && g.cfg.MicrosoftClientSecret != "" {
		providers = append(providers, microsoftonline.New(g.cfg.MicrosoftClientID, g.cfg.MicrosoftClientSecret, g.cfg.MicrosoftCallbackURL, "openid", "email", "profile", "User.Read"))
	}

	goth.UseProviders(providers...)
	return nil
}

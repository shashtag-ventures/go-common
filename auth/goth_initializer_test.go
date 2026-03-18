package auth_test

import (
	"testing"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/shashtag-ventures/go-common/auth"
	"github.com/stretchr/testify/assert"
)

func TestGothInitializer(t *testing.T) {
	cfg := auth.GothConfig{
		SessionSecret:         "test-secret",
		GoogleClientID:        "google-id",
		GoogleClientSecret:    "google-secret",
		GitHubClientID:        "github-id",
		GitHubClientSecret:    "github-secret",
		GitLabClientID:        "gitlab-id",
		GitLabClientSecret:    "gitlab-secret",
		BitbucketClientID:     "bitbucket-id",
		BitbucketClientSecret: "bitbucket-secret",
	}

	t.Run("Init success", func(t *testing.T) {
		initializer := auth.NewGothInitializer(cfg)
		err := initializer.Init()
		assert.NoError(t, err)
	})

	t.Run("Init partial config", func(t *testing.T) {
		partialCfg := auth.GothConfig{
			SessionSecret: "test-secret",
		}
		initializer := auth.NewGothInitializer(partialCfg)
		err := initializer.Init()
		assert.NoError(t, err)
	})

	t.Run("Verify cookie options", func(t *testing.T) {
		cookieCfg := auth.GothConfig{
			SessionSecret: "test-secret",
			CookieDomain:  ".example.com",
			Secure:        true,
		}
		initializer := auth.NewGothInitializer(cookieCfg)
		err := initializer.Init()
		assert.NoError(t, err)

		// Check gothic.Store (assuming it's a CookieStore)
		store, ok := gothic.Store.(*sessions.CookieStore)
		assert.True(t, ok)
		assert.Equal(t, ".example.com", store.Options.Domain)
		assert.True(t, store.Options.Secure)
		assert.Equal(t, "/", store.Options.Path)
		assert.Equal(t, 2, int(store.Options.SameSite))
	})
}

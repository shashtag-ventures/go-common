package integrations

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/crypto"
	"github.com/shashtag-ventures/go-common/testutil"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationService(t *testing.T) {
	ctx := context.Background()
	db, teardown := testutil.SetupTestDatabase(ctx)
	defer teardown()

	// Run migrations for the model
	err := db.AutoMigrate(&ExternalConnection{})
	assert.NoError(t, err)

	encryptionKey := "12345678901234567890123456789012" // 32 bytes
	storage := NewIntegrationRepository(db)
	service := NewIntegrationService(storage, encryptionKey)

	userID := uuid.New()
	provider := "github"
	accessToken := "secret-token"

	t.Run("SaveConnection encrypts tokens", func(t *testing.T) {
		err := service.SaveConnection(ctx, userID, provider, "p-user-1", "user1", "https://avatar.com", accessToken, "refresh-123")
		assert.NoError(t, err)

		// Verify in DB that it is encrypted
		var conn ExternalConnection
		err = db.First(&conn, "user_id = ?", userID).Error
		assert.NoError(t, err)
		assert.NotEqual(t, accessToken, conn.AccessToken)

		// Verify we can decrypt it back
		decrypted, err := crypto.Decrypt(conn.AccessToken, encryptionKey)
		assert.NoError(t, err)
		assert.Equal(t, accessToken, decrypted)
	})

	t.Run("GetUserConnections returns all connections for a user", func(t *testing.T) {
		conns, err := service.GetUserConnections(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, conns, 1)
		assert.Equal(t, provider, conns[0].Provider)
	})

	t.Run("ListUserRepositories decrypts and calls client", func(t *testing.T) {
		// Mock GitHub API
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "token "+accessToken, r.Header.Get("Authorization"))
			w.Write([]byte(`[{"name": "repo1"}]`))
		}))
		defer server.Close()

		impl := service.(*integrationService)
		impl.githubClient.BaseURL = server.URL

		repos, err := service.ListUserRepositories(ctx, userID, provider)
		assert.NoError(t, err)
		assert.Len(t, repos, 1)
		assert.Equal(t, "repo1", repos[0].Name)
	})

	t.Run("ListUserRepositories returns error for unsupported provider", func(t *testing.T) {
		// First save a connection for an unsupported provider
		err := service.SaveConnection(ctx, userID, "unsupported", "p-user-2", "user2", "", "token", "")
		assert.NoError(t, err)

		_, err = service.ListUserRepositories(ctx, userID, "unsupported")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not supported")
	})

	t.Run("ListUserNamespaces decrypts and calls client", func(t *testing.T) {
		// Mock GitHub API
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/user" {
				w.Write([]byte(`{"login": "user1"}`))
				return
			}
			w.Write([]byte(`{"installations": []}`))
		}))
		defer server.Close()

		impl := service.(*integrationService)
		impl.githubClient.BaseURL = server.URL

		namespaces, err := service.ListUserNamespaces(ctx, userID, provider)
		assert.NoError(t, err)
		assert.Len(t, namespaces, 1)
		assert.Equal(t, "user1", namespaces[0].Name)
	})
}

package integrations_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shashtag-ventures/go-common/integrations"
	"github.com/shashtag-ventures/go-common/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationRepository(t *testing.T) {
	ctx := context.Background()
	db, teardown := testutil.SetupTestDatabase(ctx)
	defer teardown()

	// Migrate the model
	err := db.AutoMigrate(&integrations.ExternalConnection{})
	require.NoError(t, err)

	repo := integrations.NewIntegrationRepository(db)
	userID := uuid.New()

	t.Run("Save and Get Connection", func(t *testing.T) {
		testutil.CleanTables(db, "external_connections")

		conn := &integrations.ExternalConnection{
			UserID:         userID,
			Provider:       "github",
			ProviderUserID: "github-123",
			Username:       "testuser",
		}

		err := repo.SaveConnection(ctx, conn)
		assert.NoError(t, err)

		found, err := repo.GetConnection(ctx, userID, "github")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "testuser", found.Username)
		assert.Equal(t, "github-123", found.ProviderUserID)
	})

	t.Run("Upsert Connection", func(t *testing.T) {
		testutil.CleanTables(db, "external_connections")

		conn1 := &integrations.ExternalConnection{
			UserID:   userID,
			Provider: "google",
			Username: "user1",
		}
		err := repo.SaveConnection(ctx, conn1)
		assert.NoError(t, err)

		conn2 := &integrations.ExternalConnection{
			UserID:   userID,
			Provider: "google",
			Username: "user2-updated",
		}
		err = repo.SaveConnection(ctx, conn2)
		assert.NoError(t, err)

		found, err := repo.GetConnection(ctx, userID, "google")
		assert.NoError(t, err)
		assert.Equal(t, "user2-updated", found.Username)
	})

	t.Run("List Connections", func(t *testing.T) {
		testutil.CleanTables(db, "external_connections")

		repo.SaveConnection(ctx, &integrations.ExternalConnection{
			UserID:   userID,
			Provider: "github",
		})
		repo.SaveConnection(ctx, &integrations.ExternalConnection{
			UserID:   userID,
			Provider: "google",
		})

		list, err := repo.ListConnections(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("Get Non-existent Connection", func(t *testing.T) {
		testutil.CleanTables(db, "external_connections")
		found, err := repo.GetConnection(ctx, userID, "non-existent")
		assert.Error(t, err)
		assert.Nil(t, found)
	})
}

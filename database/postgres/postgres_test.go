package postgres_test

import (
	"context"
	"testing"
	"time"

	db_pkg "github.com/shashtag-ventures/go-common/database/postgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDatabase(t *testing.T) {
	ctx := context.Background()

	// Success case
	t.Run("Successful Connection", func(t *testing.T) {
		pgContainer, err := postgres.RunContainer(ctx,
			testcontainers.WithImage("postgres:16-alpine"),
			postgres.WithDatabase("test-db"),
			postgres.WithUsername("postgres"),
			postgres.WithPassword("postgres"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second),
			),
		)
		require.NoError(t, err)
		defer pgContainer.Terminate(ctx)

		host, _ := pgContainer.Host(ctx)
		port, _ := pgContainer.MappedPort(ctx, "5432")

		cfg := db_pkg.DBConfig{
			Host:     host,
			Port:     port.Port(),
			User:     "postgres",
			Password: "postgres",
			DbName:   "test-db",
			SSLMode:  "disable",
		}

		db, err := db_pkg.New(cfg)
		assert.NoError(t, err)
		if assert.NotNil(t, db) {
			sqlDb, err := db.DB()
			assert.NoError(t, err)
			assert.NoError(t, sqlDb.Ping())
		}
	})

	// Failure case
	t.Run("Connection Failure", func(t *testing.T) {
		cfg := db_pkg.DBConfig{
			Host:     "localhost",
			Port:     "54321", // Wrong port
			User:     "invalid",
			Password: "wrong",
			DbName:   "non-existent",
			SSLMode:  "disable",
		}

		db, err := db_pkg.New(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})

	t.Run("Invalid Config Failure", func(t *testing.T) {
		cfg := db_pkg.DBConfig{
			SSLMode: "invalid-mode", // Should cause connection string error
		}
		db, err := db_pkg.New(cfg)
		assert.Error(t, err)
		assert.Nil(t, db)
	})
}

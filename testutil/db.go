package testutil

import (
	"context"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDatabase starts a PostgreSQL container and returns a GORM DB instance.
// It does NOT perform any migrations.
func SetupTestDatabase(ctx context.Context) (*gorm.DB, func()) {
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
	if err != nil {
		log.Fatalf("failed to start postgres container: %s", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	db, err := gorm.Open(gormpostgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	teardown := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate postgres container: %s", err)
		}
	}

	return db, teardown
}

// CleanTables truncates the specified tables and resets their auto-increment counters.
func CleanTables(db *gorm.DB, tables ...string) {
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE")
	}
}

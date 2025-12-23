package postgres

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/shashtag-ventures/go-common/gormutil"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	tracing "gorm.io/plugin/opentelemetry/tracing"
)

// DBConfig defines the configuration for a PostgreSQL database connection.
type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int // In minutes
}

// New initializes a new PostgreSQL database connection using GORM.
// It configures connection pooling and registers OpenTelemetry tracing.
func New(cfg DBConfig) (*gorm.DB, error) {
	// Construct the DSN (Data Source Name) for the PostgreSQL connection.
	cnn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SSLMode)

	// Open a new GORM database connection with custom contextual logger.
	dbClient, err := gorm.Open(postgres.Open(cnn), &gorm.Config{
		Logger: &gormutil.GormLogger{
			SlowThreshold: 200 * time.Millisecond,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Register OpenTelemetry plugin for GORM to enable automatic tracing of database operations.
	dbClient.Use(tracing.NewPlugin())

	// Get the underlying sql.DB object to configure connection pooling.
	sqlDb, err := dbClient.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings based on configuration.
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	// Ping the database to verify the connection is alive.
	err = sqlDb.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("Database connected successfully")

	return dbClient, nil
}
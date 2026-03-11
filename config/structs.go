package config

import "log/slog"

// GitHubConfig holds credentials for GitHub OAuth.
type GitHubConfig struct {
	ClientID     string
	ClientSecret string
}

// GitLabConfig holds credentials for GitLab OAuth.
type GitLabConfig struct {
	ClientID     string
	ClientSecret string
}

// BitbucketConfig holds credentials for Bitbucket OAuth.
type BitbucketConfig struct {
	ClientID     string
	ClientSecret string
}

// MicrosoftConfig holds credentials for Microsoft OAuth.
type MicrosoftConfig struct {
	ClientID     string
	ClientSecret string
}

// RateLimitConfig defines settings for API rate limiting.
type RateLimitConfig struct {
	Enabled bool
	Limit   int
	Window  int // Window in seconds
}

// CorsConfig defines settings for Cross-Origin Resource Sharing.
type CorsConfig struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

// GoogleConfig holds credentials for Google OAuth.
type GoogleConfig struct {
	ClientID     string
	ClientSecret string
}

// JWTConfig holds settings for JSON Web Token handling.
type JWTConfig struct {
	Secret string
}

// SessionConfig holds settings for session management.
type SessionConfig struct {
	Secret string
}

// ServerConfig defines HTTP server settings.
type ServerConfig struct {
	Addr            string
	PublicURL       string
	ShutdownTimeout int
}

// DatabaseConfig defines database connection settings.
type DatabaseConfig struct {
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

// LoggerConfig defines logging settings.
type LoggerConfig struct {
	Level slog.Level
}

// RedisConfig defines Redis connection settings.
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

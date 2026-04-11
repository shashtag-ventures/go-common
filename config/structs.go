package config

import (
	"log/slog"
	"strconv"
)

// GitHubConfig holds credentials for GitHub OAuth and GitHub App.
type GitHubConfig struct {
	ClientID     string `env:"GITHUB_CLIENT_ID"`
	ClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	AppID        string `env:"GITHUB_APP_ID"`   // GitHub App ID for installation token generation
	PrivateKey    string `env:"GITHUB_PRIVATE_KEY"` // PEM private key (used if directly in env)
	AppName       string `env:"GITHUB_APP_NAME"`    // GitHub App slug name (for installation URLs)
	WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	APIURL        string `env:"GITHUB_API_URL" envDefault:"https://api.github.com"`
}

// GitLabConfig holds credentials for GitLab OAuth.
type GitLabConfig struct {
	ClientID     string `env:"GITLAB_CLIENT_ID"`
	ClientSecret string `env:"GITLAB_CLIENT_SECRET"`
}

// BitbucketConfig holds credentials for Bitbucket OAuth.
type BitbucketConfig struct {
	ClientID     string `env:"BITBUCKET_CLIENT_ID"`
	ClientSecret string `env:"BITBUCKET_CLIENT_SECRET"`
}

// MicrosoftConfig holds credentials for Microsoft OAuth.
type MicrosoftConfig struct {
	ClientID     string `env:"MICROSOFT_CLIENT_ID"`
	ClientSecret string `env:"MICROSOFT_CLIENT_SECRET"`
}

// RateLimitConfig defines settings for API rate limiting.
type RateLimitConfig struct {
	Enabled bool `env:"RATE_LIMIT_ENABLED" envDefault:"true"`
	Limit   int  `env:"RATE_LIMIT_LIMIT" envDefault:"100"`
	Window  int  `env:"RATE_LIMIT_WINDOW_SECONDS" envDefault:"60"` // Window in seconds
}

// CorsConfig defines settings for Cross-Origin Resource Sharing.
type CorsConfig struct {
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods []string `env:"CORS_ALLOWED_METHODS"`
	AllowedHeaders []string `env:"CORS_ALLOWED_HEADERS"`
}

// GoogleConfig holds credentials for Google OAuth.
type GoogleConfig struct {
	ClientID     string `env:"GOOGLE_CLIENT_ID"`
	ClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}

// JWTConfig holds settings for JSON Web Token handling.
type JWTConfig struct {
	Secret            string `env:"JWT_SECRET"`
	ExpirationSeconds int    `env:"JWT_EXPIRATION_SECONDS" envDefault:"86400"`
}

// SessionConfig holds settings for session management.
type SessionConfig struct {
	Secret string `env:"SESSION_SECRET"`
}

// ServerConfig defines HTTP server settings.
type ServerConfig struct {
	Addr            string `env:"SERVER_PORT"`
	PublicURL       string `env:"SERVER_PUBLIC_URL"`
	ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT" envDefault:"10"`
}

// DatabaseConfig defines database connection settings.
type DatabaseConfig struct {
	Host            string `env:"DB_HOST"`
	Port            string `env:"DB_PORT"`
	User            string `env:"DB_USER"`
	Password        string `env:"DB_PASSWORD"`
	DbName          string `env:"DB_NAME"`
	SSLMode         string `env:"DB_SSL_MODE"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" envDefault:"100"`
	ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME" envDefault:"5"` // In minutes
}

// Level wraps slog.Level to provide a more robust unmarshaler that handles both numeric strings and names.
type Level slog.Level

func (l *Level) UnmarshalText(text []byte) error {
	s := string(text)
	if i, err := strconv.Atoi(s); err == nil {
		*l = Level(i)
		return nil
	}
	var sl slog.Level
	if err := sl.UnmarshalText(text); err != nil {
		return err
	}
	*l = Level(sl)
	return nil
}

// Slog returns the underlying slog.Level.
func (l Level) Slog() slog.Level {
	return slog.Level(l)
}

// LoggerConfig defines logging settings.
type LoggerConfig struct {
	Level Level `env:"LOG_LEVEL"`
}

// RedisConfig defines Redis connection settings.
type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

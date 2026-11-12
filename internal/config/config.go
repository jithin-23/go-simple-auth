package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all settings for the app, loaded once at startup.
type Config struct {
	Port    string
	GinMode string
	DB      DBConfig
	JWT     JWTConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

// DSN builds the Postgres connection string GORM needs.
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// Load reads .env (if present) and environment variables into a Config.
func Load() (*Config, error) {
	// .env is optional: in production, real environment variables are used instead.
	// Variables already set in the shell are NOT overwritten by .env.
	_ = godotenv.Load()

	ttl, err := time.ParseDuration(getEnv("JWT_TTL", "1h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL: %w", err)
	}

	cfg := &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5434"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "go_simple_auth"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
			TTL:    ttl,
		},
	}

	if cfg.JWT.Secret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	if cfg.DB.Password == "" {
		return nil, errors.New("DB_PASSWORD is required")
	}

	return cfg, nil
}

// getEnv returns the env var's value, or fallback if it's unset or empty.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

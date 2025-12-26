package config

import (
	"errors"
	"os"
)

// Config holds application configuration
type Config struct {
	DatabaseURL    string
	JWTSecret      string
	Port           string
	AllowedOrigins string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000" // Default origin
	}

	return &Config{
		DatabaseURL:    databaseURL,
		JWTSecret:      jwtSecret,
		Port:           port,
		AllowedOrigins: allowedOrigins,
	}, nil
}

package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DATABASE_URL", "postgresql://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("PORT", "9000")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
		os.Unsetenv("ALLOWED_ORIGINS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.DatabaseURL != "postgresql://test:test@localhost:5432/test" {
		t.Errorf("Expected DatabaseURL to be 'postgresql://test:test@localhost:5432/test', got '%s'", cfg.DatabaseURL)
	}

	if cfg.JWTSecret != "test-secret" {
		t.Errorf("Expected JWTSecret to be 'test-secret', got '%s'", cfg.JWTSecret)
	}

	if cfg.Port != "9000" {
		t.Errorf("Expected Port to be '9000', got '%s'", cfg.Port)
	}

	if cfg.AllowedOrigins != "http://localhost:3000,http://localhost:5173" {
		t.Errorf("Expected AllowedOrigins to be 'http://localhost:3000,http://localhost:5173', got '%s'", cfg.AllowedOrigins)
	}
}

func TestLoadMissingRequired(t *testing.T) {
	// Clear all env vars
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("PORT")
	os.Unsetenv("ALLOWED_ORIGINS")

	_, err := Load()
	if err == nil {
		t.Error("Expected Load() to fail when DATABASE_URL is missing, but it succeeded")
	}
}

func TestLoadWithDefaults(t *testing.T) {
	// Set only required env vars
	os.Setenv("DATABASE_URL", "postgresql://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
		os.Unsetenv("ALLOWED_ORIGINS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// PORT should default to 8080
	if cfg.Port != "8080" {
		t.Errorf("Expected Port to default to '8080', got '%s'", cfg.Port)
	}

	// ALLOWED_ORIGINS should default to localhost:3000
	if cfg.AllowedOrigins != "http://localhost:3000" {
		t.Errorf("Expected AllowedOrigins to default to 'http://localhost:3000', got '%s'", cfg.AllowedOrigins)
	}
}

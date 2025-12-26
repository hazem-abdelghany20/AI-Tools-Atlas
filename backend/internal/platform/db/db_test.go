package db

import (
	"testing"

	"github.com/your-org/ai-tools-atlas-backend/internal/platform/config"
)

func TestConnect(t *testing.T) {
	// Skip if no test database available
	t.Skip("Skipping database connection test - requires PostgreSQL instance")

	cfg := &config.Config{
		DatabaseURL: "postgresql://postgres:password@localhost:5432/ai_tools_atlas_test?sslmode=disable",
	}

	db, err := Connect(cfg)
	if err != nil {
		t.Fatalf("Connect() failed: %v", err)
	}

	if db == nil {
		t.Error("Expected non-nil database connection")
	}

	// Test that we can ping the database
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying *sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Errorf("Failed to ping database: %v", err)
	}

	// Clean up
	Close(db)
}

func TestClose(t *testing.T) {
	// Skip if no test database available
	t.Skip("Skipping database close test - requires PostgreSQL instance")

	cfg := &config.Config{
		DatabaseURL: "postgresql://postgres:password@localhost:5432/ai_tools_atlas_test?sslmode=disable",
	}

	db, err := Connect(cfg)
	if err != nil {
		t.Fatalf("Connect() failed: %v", err)
	}

	err = Close(db)
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}
}

func TestConnectInvalidDSN(t *testing.T) {
	cfg := &config.Config{
		DatabaseURL: "invalid-dsn",
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Error("Expected Connect() to fail with invalid DSN, but it succeeded")
	}
}

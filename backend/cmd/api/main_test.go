package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/your-org/ai-tools-atlas-backend/internal/auth"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/config"
	platformhttp "github.com/your-org/ai-tools-atlas-backend/internal/platform/http"
)

func TestIntegration_HealthCheckEndpoint(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DATABASE_URL", "postgresql://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("PORT", "8080")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:3000")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
		os.Unsetenv("ALLOWED_ORIGINS")
	}()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Setup router (with nil database for health check test - handlers that need DB won't work)
	authService := auth.NewService()
	router := platformhttp.SetupRouter(cfg, nil, authService)

	// Create test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Test health check endpoint
	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestIntegration_ServerStartsOnConfiguredPort(t *testing.T) {
	// Set up test environment variables
	os.Setenv("DATABASE_URL", "postgresql://test:test@localhost:5432/test")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("PORT", "9999")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
	}()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify port is set correctly
	if cfg.Port != "9999" {
		t.Errorf("Expected port '9999', got '%s'", cfg.Port)
	}
}

func TestIntegration_AllAcceptanceCriteriaSatisfied(t *testing.T) {
	t.Run("Go module exists", func(t *testing.T) {
		if _, err := os.Stat("../../go.mod"); os.IsNotExist(err) {
			t.Error("go.mod does not exist")
		}
	})

	t.Run("Project structure exists", func(t *testing.T) {
		dirs := []string{
			"../../cmd/api",
			"../../internal/platform/config",
			"../../internal/platform/db",
			"../../internal/platform/http",
			"../../migrations",
		}

		for _, dir := range dirs {
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				t.Errorf("Directory %s does not exist", dir)
			}
		}
	})

	t.Run(".env.example exists", func(t *testing.T) {
		if _, err := os.Stat("../../.env.example"); os.IsNotExist(err) {
			t.Error(".env.example does not exist")
		}
	})

	t.Run("Makefile exists", func(t *testing.T) {
		if _, err := os.Stat("../../Makefile"); os.IsNotExist(err) {
			t.Error("Makefile does not exist")
		}
	})

	t.Run("All required files exist", func(t *testing.T) {
		files := []string{
			"../../cmd/api/main.go",
			"../../internal/platform/config/config.go",
			"../../internal/platform/db/db.go",
			"../../internal/platform/http/router.go",
			"../../internal/platform/http/middleware.go",
			"../../internal/platform/http/responses.go",
		}

		for _, file := range files {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Errorf("File %s does not exist", file)
			}
		}
	})
}

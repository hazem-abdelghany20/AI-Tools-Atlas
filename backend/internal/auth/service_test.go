package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestMain(m *testing.M) {
	// Set JWT_SECRET for tests
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	code := m.Run()
	os.Exit(code)
}

func TestGenerateToken(t *testing.T) {
	service := NewService()

	t.Run("generates valid JWT token", func(t *testing.T) {
		userID := uint(123)
		email := "test@example.com"
		role := "user"

		tokenString, err := service.GenerateToken(userID, email, role)
		if err != nil {
			t.Fatalf("GenerateToken failed: %v", err)
		}

		if tokenString == "" {
			t.Error("Expected non-empty token string")
		}

		// Parse and verify token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return service.jwtSecret, nil
		})

		if err != nil {
			t.Fatalf("Failed to parse generated token: %v", err)
		}

		if !token.Valid {
			t.Error("Token should be valid")
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			t.Fatal("Failed to extract claims from token")
		}

		if claims.UserID != userID {
			t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
		}

		if claims.Email != email {
			t.Errorf("Expected Email %s, got %s", email, claims.Email)
		}

		if claims.Role != role {
			t.Errorf("Expected Role %s, got %s", role, claims.Role)
		}
	})

	t.Run("token has correct expiration (7 days)", func(t *testing.T) {
		tokenString, _ := service.GenerateToken(1, "test@example.com", "user")

		token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return service.jwtSecret, nil
		})

		claims := token.Claims.(*Claims)
		expectedExpiration := time.Now().Add(7 * 24 * time.Hour)
		actualExpiration := claims.ExpiresAt.Time

		// Allow 1 second tolerance
		diff := actualExpiration.Sub(expectedExpiration)
		if diff > time.Second || diff < -time.Second {
			t.Errorf("Token expiration not set to 7 days. Expected ~%v, got %v", expectedExpiration, actualExpiration)
		}
	})
}

func TestValidateToken(t *testing.T) {
	service := NewService()

	t.Run("validates correct token", func(t *testing.T) {
		// Generate a token
		userID := uint(456)
		email := "valid@example.com"
		role := "admin"
		tokenString, _ := service.GenerateToken(userID, email, role)

		// Validate it
		claims, err := service.ValidateToken(tokenString)
		if err != nil {
			t.Fatalf("ValidateToken failed for valid token: %v", err)
		}

		if claims.UserID != userID {
			t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
		}

		if claims.Email != email {
			t.Errorf("Expected Email %s, got %s", email, claims.Email)
		}

		if claims.Role != role {
			t.Errorf("Expected Role %s, got %s", role, claims.Role)
		}
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		_, err := service.ValidateToken(invalidToken)
		if err == nil {
			t.Error("Expected error for invalid token, got nil")
		}
	})

	t.Run("rejects token with wrong signature", func(t *testing.T) {
		// Create token with different secret
		wrongService := &Service{
			jwtSecret:     []byte("wrong-secret"),
			tokenDuration: 7 * 24 * time.Hour,
		}
		tokenString, _ := wrongService.GenerateToken(1, "test@example.com", "user")

		// Try to validate with correct service
		_, err := service.ValidateToken(tokenString)
		if err == nil {
			t.Error("Expected error for token with wrong signature, got nil")
		}
	})

	t.Run("rejects expired token", func(t *testing.T) {
		// Create service with very short expiration
		shortLivedService := &Service{
			jwtSecret:     service.jwtSecret,
			tokenDuration: -1 * time.Hour, // Already expired
		}

		tokenString, _ := shortLivedService.GenerateToken(1, "test@example.com", "user")

		_, err := service.ValidateToken(tokenString)
		if err == nil {
			t.Error("Expected error for expired token, got nil")
		}
	})
}

func TestHashPassword(t *testing.T) {
	service := NewService()

	t.Run("creates bcrypt hash", func(t *testing.T) {
		password := "mySecurePassword123"

		hash, err := service.HashPassword(password)
		if err != nil {
			t.Fatalf("HashPassword failed: %v", err)
		}

		if hash == "" {
			t.Error("Expected non-empty hash")
		}

		if hash == password {
			t.Error("Hash should not equal plain password")
		}

		// Hash should start with bcrypt prefix
		if len(hash) < 10 {
			t.Error("Hash seems too short for bcrypt")
		}
	})

	t.Run("generates different hashes for same password", func(t *testing.T) {
		password := "samePassword"

		hash1, _ := service.HashPassword(password)
		hash2, _ := service.HashPassword(password)

		// Due to bcrypt salt, hashes should be different
		if hash1 == hash2 {
			t.Error("Expected different hashes for same password (bcrypt uses random salt)")
		}
	})
}

func TestCheckPassword(t *testing.T) {
	service := NewService()

	t.Run("validates correct password", func(t *testing.T) {
		password := "correctPassword123"
		hash, _ := service.HashPassword(password)

		err := service.CheckPassword(password, hash)
		if err != nil {
			t.Errorf("CheckPassword failed for correct password: %v", err)
		}
	})

	t.Run("rejects wrong password", func(t *testing.T) {
		password := "correctPassword"
		hash, _ := service.HashPassword(password)

		err := service.CheckPassword("wrongPassword", hash)
		if err == nil {
			t.Error("Expected error for wrong password, got nil")
		}
	})

	t.Run("rejects empty password", func(t *testing.T) {
		password := "password"
		hash, _ := service.HashPassword(password)

		err := service.CheckPassword("", hash)
		if err == nil {
			t.Error("Expected error for empty password, got nil")
		}
	})

	t.Run("rejects invalid hash", func(t *testing.T) {
		password := "password"
		invalidHash := "not-a-valid-bcrypt-hash"

		err := service.CheckPassword(password, invalidHash)
		if err == nil {
			t.Error("Expected error for invalid hash, got nil")
		}
	})
}

func TestNewService(t *testing.T) {
	t.Run("creates service with JWT_SECRET from env", func(t *testing.T) {
		service := NewService()

		if service == nil {
			t.Fatal("Expected non-nil service")
		}

		if len(service.jwtSecret) == 0 {
			t.Error("Service should have jwtSecret loaded from env")
		}

		expectedDuration := 7 * 24 * time.Hour
		if service.tokenDuration != expectedDuration {
			t.Errorf("Expected token duration %v, got %v", expectedDuration, service.tokenDuration)
		}
	})

	t.Run("panics if JWT_SECRET not set", func(t *testing.T) {
		// Save and unset JWT_SECRET
		originalSecret := os.Getenv("JWT_SECRET")
		os.Unsetenv("JWT_SECRET")

		defer func() {
			// Restore JWT_SECRET
			os.Setenv("JWT_SECRET", originalSecret)

			// Verify panic occurred
			if r := recover(); r == nil {
				t.Error("Expected panic when JWT_SECRET not set, got nil")
			}
		}()

		// Should panic
		NewService()
	})
}

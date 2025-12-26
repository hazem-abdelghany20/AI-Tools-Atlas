package http

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/auth"
)

func TestMain(m *testing.M) {
	// Set JWT_SECRET for tests
	os.Setenv("JWT_SECRET", "test-secret-for-middleware-tests")
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func TestAuthRequired(t *testing.T) {
	authService := auth.NewService()

	t.Run("allows valid token", func(t *testing.T) {
		// Generate valid token
		tokenString, _ := authService.GenerateToken(123, "test@example.com", "user")

		// Setup test router
		router := gin.New()
		router.Use(AuthRequired(authService))
		router.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			userRole, _ := c.Get("user_role")
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"role":    userRole,
			})
		})

		// Make request with cookie
		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: tokenString,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("blocks missing token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthRequired(authService))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		// No cookie set
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("blocks invalid token", func(t *testing.T) {
		router := gin.New()
		router.Use(AuthRequired(authService))
		router.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: "invalid.token.value",
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status 401, got %d", w.Code)
		}
	})

	t.Run("sets user context correctly", func(t *testing.T) {
		userID := uint(789)
		email := "context@example.com"
		role := "moderator"
		tokenString, _ := authService.GenerateToken(userID, email, role)

		router := gin.New()
		router.Use(AuthRequired(authService))
		router.GET("/protected", func(c *gin.Context) {
			contextUserID, exists1 := c.Get("user_id")
			contextRole, exists2 := c.Get("user_role")
			contextEmail, exists3 := c.Get("user_email")

			if !exists1 || !exists2 || !exists3 {
				t.Error("Context values not set")
			}

			if contextUserID != userID {
				t.Errorf("Expected user_id %d, got %v", userID, contextUserID)
			}

			if contextRole != role {
				t.Errorf("Expected role %s, got %v", role, contextRole)
			}

			if contextEmail != email {
				t.Errorf("Expected email %s, got %v", email, contextEmail)
			}

			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: tokenString,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})
}

func TestAdminRequired(t *testing.T) {
	t.Run("allows admin role", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			// Simulate authenticated user with admin role
			c.Set("user_role", "admin")
			c.Next()
		})
		router.Use(AdminRequired())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("blocks non-admin role", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			// Simulate authenticated user with user role
			c.Set("user_role", "user")
			c.Next()
		})
		router.Use(AdminRequired())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})

	t.Run("blocks moderator role", func(t *testing.T) {
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("user_role", "moderator")
			c.Next()
		})
		router.Use(AdminRequired())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})

	t.Run("blocks when role not set in context", func(t *testing.T) {
		router := gin.New()
		// No role set in context
		router.Use(AdminRequired())
		router.GET("/admin", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		})

		req := httptest.NewRequest("GET", "/admin", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})
}

func TestAuthRequiredAndAdminRequiredChain(t *testing.T) {
	authService := auth.NewService()

	t.Run("allows admin user through both middlewares", func(t *testing.T) {
		// Generate admin token
		tokenString, _ := authService.GenerateToken(1, "admin@example.com", "admin")

		router := gin.New()
		router.Use(AuthRequired(authService))
		router.Use(AdminRequired())
		router.GET("/admin-only", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin only content"})
		})

		req := httptest.NewRequest("GET", "/admin-only", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: tokenString,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("blocks non-admin user", func(t *testing.T) {
		// Generate user token (not admin)
		tokenString, _ := authService.GenerateToken(2, "user@example.com", "user")

		router := gin.New()
		router.Use(AuthRequired(authService))
		router.Use(AdminRequired())
		router.GET("/admin-only", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "admin only content"})
		})

		req := httptest.NewRequest("GET", "/admin-only", nil)
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: tokenString,
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403, got %d", w.Code)
		}
	})
}

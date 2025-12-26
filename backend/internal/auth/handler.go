package auth

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/bookmarks"
)

// Handler handles authentication HTTP requests
type Handler struct {
	service         *Service
	bookmarkService bookmarks.Service
}

// NewHandler creates a new auth handler
func NewHandler(service *Service, bookmarkService bookmarks.Service) *Handler {
	return &Handler{
		service:         service,
		bookmarkService: bookmarkService,
	}
}

// RegisterRoutes sets up the auth routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
	}

	// Protected route for current user
	rg.GET("/me", authMiddleware, h.GetCurrentUser)
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request body",
			},
		})
		return
	}

	user, token, err := h.service.Register(input)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidEmail):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{
					"code":    "INVALID_EMAIL",
					"message": err.Error(),
				},
			})
		case errors.Is(err, ErrPasswordTooShort):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{
					"code":    "PASSWORD_TOO_SHORT",
					"message": err.Error(),
				},
			})
		case errors.Is(err, ErrDisplayNameRequired):
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": gin.H{
					"code":    "DISPLAY_NAME_REQUIRED",
					"message": err.Error(),
				},
			})
		case errors.Is(err, ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "EMAIL_EXISTS",
					"message": "An account with this email already exists",
				},
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "Failed to create account",
				},
			})
		}
		return
	}

	// Migrate session bookmarks to the new user
	h.migrateSessionBookmarks(c, user.ID)

	// Set auth cookie
	h.setAuthCookie(c, token)

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"user": ToUserResponse(user),
		},
	})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request body",
			},
		})
		return
	}

	user, token, err := h.service.Login(input)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_CREDENTIALS",
					"message": "Invalid email or password",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to login",
			},
		})
		return
	}

	// Migrate session bookmarks to the logged-in user
	h.migrateSessionBookmarks(c, user.ID)

	// Set auth cookie
	h.setAuthCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"user": ToUserResponse(user),
		},
	})
}

// Logout handles user logout
func (h *Handler) Logout(c *gin.Context) {
	// Clear the auth cookie
	h.clearAuthCookie(c)
	c.Status(http.StatusNoContent)
}

// GetCurrentUser returns the current authenticated user
func (h *Handler) GetCurrentUser(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Authentication required",
			},
		})
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Invalid user context",
			},
		})
		return
	}

	user, err := h.service.GetCurrentUser(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "USER_NOT_FOUND",
					"message": "User not found",
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Failed to get user",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ToUserResponse(user),
	})
}

// setAuthCookie sets the JWT token as an HTTP-only cookie
func (h *Handler) setAuthCookie(c *gin.Context, token string) {
	secure := os.Getenv("APP_ENV") == "production"
	maxAge := 7 * 24 * 60 * 60 // 7 days in seconds

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"auth_token",
		token,
		maxAge,
		"/",
		"",     // domain
		secure, // secure
		true,   // httpOnly
	)
}

// clearAuthCookie removes the auth cookie
func (h *Handler) clearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"auth_token",
		"",
		-1,    // MaxAge < 0 means delete
		"/",
		"",    // domain
		false, // secure
		true,  // httpOnly
	)
}

// migrateSessionBookmarks moves anonymous session bookmarks to the user
func (h *Handler) migrateSessionBookmarks(c *gin.Context, userID uint) {
	sessionID, err := c.Cookie("session_id")
	if err != nil || sessionID == "" {
		return
	}

	if h.bookmarkService != nil {
		// Migrate bookmarks (ignore errors as this is optional)
		_ = h.bookmarkService.MigrateSessionBookmarks(userID, sessionID)
	}

	// Clear the session cookie after migration
	c.SetCookie(
		"session_id",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}

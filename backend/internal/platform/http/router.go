package http

import (
	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/analytics"
	"github.com/your-org/ai-tools-atlas-backend/internal/auth"
	"github.com/your-org/ai-tools-atlas-backend/internal/badges"
	"github.com/your-org/ai-tools-atlas-backend/internal/bookmarks"
	"github.com/your-org/ai-tools-atlas-backend/internal/categories"
	"github.com/your-org/ai-tools-atlas-backend/internal/moderation"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/config"
	"github.com/your-org/ai-tools-atlas-backend/internal/reviews"
	"github.com/your-org/ai-tools-atlas-backend/internal/tags"
	"github.com/your-org/ai-tools-atlas-backend/internal/tools"
	"gorm.io/gorm"
)

// SetupRouter initializes the Gin router with middleware
func SetupRouter(cfg *config.Config, db *gorm.DB, authService *auth.Service) *gin.Engine {
	// Create router
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware(cfg.AllowedOrigins))
	r.Use(RequestLoggerMiddleware())

	// Health check endpoint (outside versioned API)
	r.GET("/health", HealthCheck)

	// Create API v1 group
	v1 := r.Group("/api/v1")

	// Auth middleware for protected routes
	authMiddleware := AuthRequired(authService)
	optionalAuthMiddleware := OptionalAuth(authService)
	adminMiddleware := AdminRequired()

	// Initialize repositories
	authRepo := auth.NewRepository(db)
	categoryRepo := categories.NewRepository(db)
	toolRepo := tools.NewRepository(db)
	reviewRepo := reviews.NewRepository(db)
	bookmarkRepo := bookmarks.NewRepository(db)
	tagRepo := tags.NewRepository(db)
	badgeRepo := badges.NewRepository(db)
	analyticsRepo := analytics.NewRepository(db)
	moderationRepo := moderation.NewRepository(db)

	// Initialize services
	// Update auth service with repository for register/login
	authServiceWithRepo := auth.NewServiceWithRepo(authRepo)
	categoryService := categories.NewService(categoryRepo)
	toolService := tools.NewService(toolRepo)
	reviewService := reviews.NewService(reviewRepo)
	bookmarkService := bookmarks.NewService(bookmarkRepo)
	tagService := tags.NewService(tagRepo)
	badgeService := badges.NewService(badgeRepo)
	analyticsService := analytics.NewService(analyticsRepo)
	moderationService := moderation.NewService(moderationRepo, reviewRepo)

	// Initialize handlers and register routes
	// Auth handler (with bookmark service for session migration)
	authHandler := auth.NewHandler(authServiceWithRepo, bookmarkService)
	authHandler.RegisterRoutes(v1, authMiddleware)

	categoryHandler := categories.NewHandler(categoryService)
	categoryHandler.RegisterRoutes(v1)

	toolHandler := tools.NewHandler(toolService)
	toolHandler.RegisterRoutes(v1)

	reviewHandler := reviews.NewHandler(reviewService)
	reviewHandler.RegisterRoutes(v1, authMiddleware)

	bookmarkHandler := bookmarks.NewHandler(bookmarkService)
	bookmarkHandler.RegisterRoutes(v1, authMiddleware, optionalAuthMiddleware)

	tagHandler := tags.NewHandler(tagService)
	badgeHandler := badges.NewHandler(badgeService)
	analyticsHandler := analytics.NewHandler(analyticsService)
	moderationHandler := moderation.NewHandler(moderationService, toolService)

	// Register moderation (reporting) public routes
	moderationHandler.RegisterRoutes(v1, optionalAuthMiddleware)

	// Admin routes (require authentication + admin role)
	admin := v1.Group("/admin")
	admin.Use(authMiddleware, adminMiddleware)
	{
		toolHandler.RegisterAdminRoutes(admin)
		categoryHandler.RegisterAdminRoutes(admin)
		tagHandler.RegisterAdminRoutes(admin)
		badgeHandler.RegisterAdminRoutes(admin)
		analyticsHandler.RegisterAdminRoutes(admin)
		moderationHandler.RegisterAdminRoutes(admin)
	}

	return r
}

package analytics

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/responses"
)

// Handler handles HTTP requests for analytics
type Handler struct {
	service Service
}

// NewHandler creates a new analytics handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterAdminRoutes registers admin analytics routes
func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	analytics := rg.Group("/analytics")
	{
		analytics.GET("/overview", h.GetOverview)
		analytics.GET("/top-tools", h.GetTopTools)
		analytics.GET("/top-categories", h.GetTopCategories)
	}
}

// GetOverview handles GET /api/v1/admin/analytics/overview
func (h *Handler) GetOverview(c *gin.Context) {
	stats, err := h.service.GetOverviewStats()
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch analytics", nil)
		return
	}
	responses.Success(c, stats)
}

// GetTopTools handles GET /api/v1/admin/analytics/top-tools
func (h *Handler) GetTopTools(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	topTools, err := h.service.GetTopTools(limit)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch top tools", nil)
		return
	}
	responses.Success(c, topTools)
}

// GetTopCategories handles GET /api/v1/admin/analytics/top-categories
func (h *Handler) GetTopCategories(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	categories, err := h.service.GetTopCategories(limit)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch top categories", nil)
		return
	}
	responses.Success(c, categories)
}

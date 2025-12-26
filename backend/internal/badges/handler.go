package badges

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/responses"
)

// Handler handles HTTP requests for badges
type Handler struct {
	service Service
}

// NewHandler creates a new badge handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterAdminRoutes registers admin badge routes on tools
func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	// Badges CRUD
	badges := rg.Group("/badges")
	{
		badges.GET("", h.ListBadges)
	}

	// Tool badges assignment
	tools := rg.Group("/tools")
	{
		tools.GET("/:id/badges", h.GetToolBadges)
		tools.POST("/:id/badges", h.AssignBadge)
		tools.DELETE("/:id/badges/:badge_id", h.RemoveBadge)
	}
}

// ListBadges handles GET /api/v1/admin/badges
func (h *Handler) ListBadges(c *gin.Context) {
	badges, err := h.service.ListBadges()
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch badges", nil)
		return
	}
	responses.Success(c, badges)
}

// GetToolBadges handles GET /api/v1/admin/tools/:id/badges
func (h *Handler) GetToolBadges(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tool ID", nil)
		return
	}

	badges, err := h.service.GetToolBadges(uint(id))
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tool badges", nil)
		return
	}
	responses.Success(c, badges)
}

// AssignBadgeInput represents input for assigning a badge
type AssignBadgeInput struct {
	BadgeID uint `json:"badge_id" binding:"required"`
}

// AssignBadge handles POST /api/v1/admin/tools/:id/badges
func (h *Handler) AssignBadge(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tool ID", nil)
		return
	}

	var input AssignBadgeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	err = h.service.AssignBadgeToTool(uint(id), input.BadgeID)
	if err != nil {
		switch {
		case errors.Is(err, ErrBadgeNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Badge not found", nil)
		case errors.Is(err, ErrBadgeAlreadyAssigned):
			responses.Error(c, http.StatusConflict, "ALREADY_ASSIGNED", "Badge already assigned to tool", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to assign badge", nil)
		}
		return
	}

	responses.Success(c, map[string]string{"message": "Badge assigned successfully"})
}

// RemoveBadge handles DELETE /api/v1/admin/tools/:id/badges/:badge_id
func (h *Handler) RemoveBadge(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tool ID", nil)
		return
	}

	badgeIDParam := c.Param("badge_id")
	badgeID, err := strconv.ParseUint(badgeIDParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid badge ID", nil)
		return
	}

	err = h.service.RemoveBadgeFromTool(uint(id), uint(badgeID))
	if err != nil {
		if errors.Is(err, ErrBadgeNotAssigned) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Badge not assigned to tool", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to remove badge", nil)
		return
	}

	c.Status(http.StatusNoContent)
}

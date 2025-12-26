package tags

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/responses"
)

// Handler handles HTTP requests for tags
type Handler struct {
	service Service
}

// NewHandler creates a new tag handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterAdminRoutes registers admin tag routes
func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	tags := rg.Group("/tags")
	{
		tags.GET("", h.AdminListTags)
		tags.POST("", h.AdminCreateTag)
		tags.GET("/:id", h.AdminGetTag)
		tags.PATCH("/:id", h.AdminUpdateTag)
		tags.DELETE("/:id", h.AdminDeleteTag)
	}
}

// AdminListTags handles GET /api/v1/admin/tags
func (h *Handler) AdminListTags(c *gin.Context) {
	tags, err := h.service.ListTagsWithCount()
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tags", nil)
		return
	}
	responses.Success(c, tags)
}

// AdminGetTag handles GET /api/v1/admin/tags/:id
func (h *Handler) AdminGetTag(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tag ID", nil)
		return
	}

	tag, err := h.service.GetTagByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrTagNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tag not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tag", nil)
		return
	}
	responses.Success(c, tag)
}

// AdminCreateTag handles POST /api/v1/admin/tags
func (h *Handler) AdminCreateTag(c *gin.Context) {
	var input CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	tag, err := h.service.CreateTag(input)
	if err != nil {
		switch {
		case errors.Is(err, ErrSlugRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Slug is required", map[string]string{"slug": "required"})
		case errors.Is(err, ErrNameRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Name is required", map[string]string{"name": "required"})
		case errors.Is(err, ErrSlugExists):
			responses.Error(c, http.StatusConflict, "SLUG_EXISTS", "A tag with this slug already exists", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create tag", nil)
		}
		return
	}
	responses.Created(c, tag)
}

// AdminUpdateTag handles PATCH /api/v1/admin/tags/:id
func (h *Handler) AdminUpdateTag(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tag ID", nil)
		return
	}

	var input UpdateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	tag, err := h.service.UpdateTag(uint(id), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrTagNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tag not found", nil)
		case errors.Is(err, ErrNameRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Name cannot be empty", map[string]string{"name": "required"})
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update tag", nil)
		}
		return
	}
	responses.Success(c, tag)
}

// AdminDeleteTag handles DELETE /api/v1/admin/tags/:id
func (h *Handler) AdminDeleteTag(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid tag ID", nil)
		return
	}

	err = h.service.DeleteTag(uint(id))
	if err != nil {
		if errors.Is(err, ErrTagNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Tag not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete tag", nil)
		return
	}
	c.Status(http.StatusNoContent)
}

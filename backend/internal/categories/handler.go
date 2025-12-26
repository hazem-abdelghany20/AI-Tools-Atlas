package categories

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/ai-tools-atlas-backend/internal/platform/responses"
)

// Handler handles HTTP requests for categories
type Handler struct {
	service Service
}

// NewHandler creates a new category handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers category routes on the given router group
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")
	{
		categories.GET("", h.ListCategories)
		categories.GET("/:slug/tools", h.ListToolsByCategory)
	}
}

// RegisterAdminRoutes registers admin category routes
func (h *Handler) RegisterAdminRoutes(rg *gin.RouterGroup) {
	categories := rg.Group("/categories")
	{
		categories.GET("", h.AdminListCategories)
		categories.POST("", h.AdminCreateCategory)
		categories.GET("/:id", h.AdminGetCategory)
		categories.PATCH("/:id", h.AdminUpdateCategory)
		categories.DELETE("/:id", h.AdminDeleteCategory)
	}
}

// ListCategories handles GET /api/v1/categories
func (h *Handler) ListCategories(c *gin.Context) {
	cats, err := h.service.ListCategories()
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch categories", nil)
		return
	}

	responses.Success(c, cats)
}

// ListToolsByCategory handles GET /api/v1/categories/:slug/tools
func (h *Handler) ListToolsByCategory(c *gin.Context) {
	slug := c.Param("slug")

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tools, total, err := h.service.ListToolsByCategory(slug, page, pageSize)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Category not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch tools", nil)
		return
	}

	responses.List(c, tools, map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// AdminListCategories handles GET /api/v1/admin/categories
func (h *Handler) AdminListCategories(c *gin.Context) {
	cats, err := h.service.ListCategoriesWithCount()
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch categories", nil)
		return
	}
	responses.Success(c, cats)
}

// AdminGetCategory handles GET /api/v1/admin/categories/:id
func (h *Handler) AdminGetCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid category ID", nil)
		return
	}

	cat, err := h.service.GetCategoryByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Category not found", nil)
			return
		}
		responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch category", nil)
		return
	}
	responses.Success(c, cat)
}

// AdminCreateCategory handles POST /api/v1/admin/categories
func (h *Handler) AdminCreateCategory(c *gin.Context) {
	var input CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	cat, err := h.service.CreateCategory(input)
	if err != nil {
		switch {
		case errors.Is(err, ErrSlugRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Slug is required", map[string]string{"slug": "required"})
		case errors.Is(err, ErrNameRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Name is required", map[string]string{"name": "required"})
		case errors.Is(err, ErrSlugExists):
			responses.Error(c, http.StatusConflict, "SLUG_EXISTS", "A category with this slug already exists", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create category", nil)
		}
		return
	}
	responses.Created(c, cat)
}

// AdminUpdateCategory handles PATCH /api/v1/admin/categories/:id
func (h *Handler) AdminUpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid category ID", nil)
		return
	}

	var input UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", nil)
		return
	}

	cat, err := h.service.UpdateCategory(uint(id), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Category not found", nil)
		case errors.Is(err, ErrNameRequired):
			responses.Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Name cannot be empty", map[string]string{"name": "required"})
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update category", nil)
		}
		return
	}
	responses.Success(c, cat)
}

// AdminDeleteCategory handles DELETE /api/v1/admin/categories/:id
func (h *Handler) AdminDeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "INVALID_ID", "Invalid category ID", nil)
		return
	}

	err = h.service.DeleteCategory(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, ErrCategoryNotFound):
			responses.Error(c, http.StatusNotFound, "NOT_FOUND", "Category not found", nil)
		case errors.Is(err, ErrHasTools):
			responses.Error(c, http.StatusConflict, "HAS_TOOLS", "Cannot delete category with tools. Reassign tools first.", nil)
		default:
			responses.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete category", nil)
		}
		return
	}
	c.Status(http.StatusNoContent)
}

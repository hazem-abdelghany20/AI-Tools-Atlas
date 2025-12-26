package categories

import (
	"errors"

	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrSlugRequired     = errors.New("slug is required")
	ErrNameRequired     = errors.New("name is required")
	ErrSlugExists       = errors.New("slug already exists")
	ErrHasTools         = errors.New("category has tools and cannot be deleted")
)

// CreateCategoryInput represents input for creating a category
type CreateCategoryInput struct {
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	DisplayOrder int    `json:"display_order"`
}

// UpdateCategoryInput represents input for updating a category
type UpdateCategoryInput struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
}

// CategoryWithCount includes tool count
type CategoryWithCount struct {
	Category
	ToolCount int64 `json:"tool_count"`
}

// Service defines the interface for category business logic
type Service interface {
	ListCategories() ([]Category, error)
	GetCategoryBySlug(slug string) (*Category, error)
	ListToolsByCategory(slug string, page, pageSize int) ([]domain.Tool, int64, error)
	// Admin methods
	ListCategoriesWithCount() ([]CategoryWithCount, error)
	GetCategoryByID(id uint) (*Category, error)
	CreateCategory(input CreateCategoryInput) (*Category, error)
	UpdateCategory(id uint, input UpdateCategoryInput) (*Category, error)
	DeleteCategory(id uint) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new category service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ListCategories returns all active categories
func (s *service) ListCategories() ([]Category, error) {
	return s.repo.ListCategories()
}

// GetCategoryBySlug finds a category by its slug
func (s *service) GetCategoryBySlug(slug string) (*Category, error) {
	category, err := s.repo.GetCategoryBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

// ListToolsByCategory returns paginated tools for a category slug
func (s *service) ListToolsByCategory(slug string, page, pageSize int) ([]domain.Tool, int64, error) {
	// First get the category to find its ID
	category, err := s.GetCategoryBySlug(slug)
	if err != nil {
		return nil, 0, err
	}

	// Apply default pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return s.repo.ListToolsByCategory(category.ID, page, pageSize)
}

// ListCategoriesWithCount returns all categories with their tool counts
func (s *service) ListCategoriesWithCount() ([]CategoryWithCount, error) {
	cats, err := s.repo.ListCategories()
	if err != nil {
		return nil, err
	}

	result := make([]CategoryWithCount, len(cats))
	for i, cat := range cats {
		count, _ := s.repo.GetToolCount(cat.ID)
		result[i] = CategoryWithCount{
			Category:  cat,
			ToolCount: count,
		}
	}
	return result, nil
}

// GetCategoryByID finds a category by ID
func (s *service) GetCategoryByID(id uint) (*Category, error) {
	cat, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return cat, nil
}

// CreateCategory creates a new category
func (s *service) CreateCategory(input CreateCategoryInput) (*Category, error) {
	if input.Slug == "" {
		return nil, ErrSlugRequired
	}
	if input.Name == "" {
		return nil, ErrNameRequired
	}

	exists, err := s.repo.SlugExists(input.Slug, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugExists
	}

	cat := &Category{
		Slug:         input.Slug,
		Name:         input.Name,
		Description:  input.Description,
		IconURL:      input.IconURL,
		DisplayOrder: input.DisplayOrder,
	}

	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}

	return cat, nil
}

// UpdateCategory updates an existing category
func (s *service) UpdateCategory(id uint, input UpdateCategoryInput) (*Category, error) {
	cat, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	if input.Name != nil {
		if *input.Name == "" {
			return nil, ErrNameRequired
		}
		cat.Name = *input.Name
	}
	if input.Description != nil {
		cat.Description = *input.Description
	}
	if input.IconURL != nil {
		cat.IconURL = *input.IconURL
	}
	if input.DisplayOrder != nil {
		cat.DisplayOrder = *input.DisplayOrder
	}

	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}

	return cat, nil
}

// DeleteCategory deletes a category if it has no tools
func (s *service) DeleteCategory(id uint) error {
	_, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}

	count, err := s.repo.GetToolCount(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrHasTools
	}

	return s.repo.Delete(id)
}

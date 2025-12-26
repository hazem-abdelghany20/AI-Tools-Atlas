package categories

import (
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

// Repository defines the interface for category data operations
type Repository interface {
	ListCategories() ([]Category, error)
	GetCategoryBySlug(slug string) (*Category, error)
	ListToolsByCategory(categoryID uint, page, pageSize int) ([]domain.Tool, int64, error)
	// Admin methods
	GetCategoryByID(id uint) (*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id uint) error
	SlugExists(slug string, excludeID uint) (bool, error)
	GetToolCount(categoryID uint) (int64, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new category repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ListCategories returns all active categories ordered by display_order
func (r *repository) ListCategories() ([]Category, error) {
	var categories []Category
	err := r.db.Order("display_order ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryBySlug finds a category by its slug
func (r *repository) GetCategoryBySlug(slug string) (*Category, error) {
	var category Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// ListToolsByCategory returns paginated tools for a category
func (r *repository) ListToolsByCategory(categoryID uint, page, pageSize int) ([]domain.Tool, int64, error) {
	var toolsList []domain.Tool
	var total int64

	// Count total tools in category (excluding archived)
	err := r.db.Model(&domain.Tool{}).
		Where("primary_category_id = ? AND archived_at IS NULL", categoryID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated tools with preloaded relationships
	offset := (page - 1) * pageSize
	err = r.db.
		Preload("Tags").
		Preload("PrimaryCategory").
		Where("primary_category_id = ? AND archived_at IS NULL", categoryID).
		Limit(pageSize).
		Offset(offset).
		Find(&toolsList).Error
	if err != nil {
		return nil, 0, err
	}

	return toolsList, total, nil
}

// GetCategoryByID finds a category by ID
func (r *repository) GetCategoryByID(id uint) (*Category, error) {
	var category Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Create inserts a new category
func (r *repository) Create(category *Category) error {
	return r.db.Create(category).Error
}

// Update updates an existing category
func (r *repository) Update(category *Category) error {
	return r.db.Save(category).Error
}

// Delete removes a category
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Category{}, id).Error
}

// SlugExists checks if a slug is already in use
func (r *repository) SlugExists(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&Category{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetToolCount returns the number of tools in a category
func (r *repository) GetToolCount(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Tool{}).
		Where("primary_category_id = ? AND archived_at IS NULL", categoryID).
		Count(&count).Error
	return count, err
}

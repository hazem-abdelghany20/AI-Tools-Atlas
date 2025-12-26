package tags

import (
	"gorm.io/gorm"
)

// Repository defines the interface for tag data operations
type Repository interface {
	ListTags() ([]Tag, error)
	GetTagByID(id uint) (*Tag, error)
	Create(tag *Tag) error
	Update(tag *Tag) error
	Delete(id uint) error
	SlugExists(slug string, excludeID uint) (bool, error)
	GetToolCount(tagID uint) (int64, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new tag repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ListTags returns all tags
func (r *repository) ListTags() ([]Tag, error) {
	var tags []Tag
	err := r.db.Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetTagByID finds a tag by ID
func (r *repository) GetTagByID(id uint) (*Tag, error) {
	var tag Tag
	err := r.db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Create inserts a new tag
func (r *repository) Create(tag *Tag) error {
	return r.db.Create(tag).Error
}

// Update updates an existing tag
func (r *repository) Update(tag *Tag) error {
	return r.db.Save(tag).Error
}

// Delete removes a tag (will cascade via tool_tags join table)
func (r *repository) Delete(id uint) error {
	// First remove from join table
	if err := r.db.Exec("DELETE FROM tool_tags WHERE tag_id = ?", id).Error; err != nil {
		return err
	}
	return r.db.Delete(&Tag{}, id).Error
}

// SlugExists checks if a slug is already in use
func (r *repository) SlugExists(slug string, excludeID uint) (bool, error) {
	var count int64
	query := r.db.Model(&Tag{}).Where("slug = ?", slug)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetToolCount returns the number of tools with this tag
func (r *repository) GetToolCount(tagID uint) (int64, error) {
	var count int64
	err := r.db.Raw("SELECT COUNT(*) FROM tool_tags WHERE tag_id = ?", tagID).Scan(&count).Error
	return count, err
}

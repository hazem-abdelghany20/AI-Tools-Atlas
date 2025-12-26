package tags

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrTagNotFound  = errors.New("tag not found")
	ErrSlugRequired = errors.New("slug is required")
	ErrNameRequired = errors.New("name is required")
	ErrSlugExists   = errors.New("slug already exists")
)

// CreateTagInput represents input for creating a tag
type CreateTagInput struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// UpdateTagInput represents input for updating a tag
type UpdateTagInput struct {
	Name *string `json:"name,omitempty"`
}

// TagWithCount includes tool count
type TagWithCount struct {
	Tag
	ToolCount int64 `json:"tool_count"`
}

// Service defines the interface for tag business logic
type Service interface {
	ListTagsWithCount() ([]TagWithCount, error)
	GetTagByID(id uint) (*Tag, error)
	CreateTag(input CreateTagInput) (*Tag, error)
	UpdateTag(id uint, input UpdateTagInput) (*Tag, error)
	DeleteTag(id uint) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new tag service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ListTagsWithCount returns all tags with their tool counts
func (s *service) ListTagsWithCount() ([]TagWithCount, error) {
	tags, err := s.repo.ListTags()
	if err != nil {
		return nil, err
	}

	result := make([]TagWithCount, len(tags))
	for i, tag := range tags {
		count, _ := s.repo.GetToolCount(tag.ID)
		result[i] = TagWithCount{
			Tag:       tag,
			ToolCount: count,
		}
	}
	return result, nil
}

// GetTagByID finds a tag by ID
func (s *service) GetTagByID(id uint) (*Tag, error) {
	tag, err := s.repo.GetTagByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}
	return tag, nil
}

// CreateTag creates a new tag
func (s *service) CreateTag(input CreateTagInput) (*Tag, error) {
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

	tag := &Tag{
		Slug: input.Slug,
		Name: input.Name,
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// UpdateTag updates an existing tag
func (s *service) UpdateTag(id uint, input UpdateTagInput) (*Tag, error) {
	tag, err := s.repo.GetTagByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTagNotFound
		}
		return nil, err
	}

	if input.Name != nil {
		if *input.Name == "" {
			return nil, ErrNameRequired
		}
		tag.Name = *input.Name
	}

	if err := s.repo.Update(tag); err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag deletes a tag (removes from all tools)
func (s *service) DeleteTag(id uint) error {
	_, err := s.repo.GetTagByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTagNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}

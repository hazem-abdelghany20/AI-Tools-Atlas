package badges

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrBadgeNotFound       = errors.New("badge not found")
	ErrBadgeAlreadyAssigned = errors.New("badge already assigned to tool")
	ErrBadgeNotAssigned    = errors.New("badge not assigned to tool")
)

// Service defines the interface for badge business logic
type Service interface {
	ListBadges() ([]Badge, error)
	GetBadgeByID(id uint) (*Badge, error)
	AssignBadgeToTool(toolID, badgeID uint) error
	RemoveBadgeFromTool(toolID, badgeID uint) error
	GetToolBadges(toolID uint) ([]Badge, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new badge service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ListBadges returns all badges
func (s *service) ListBadges() ([]Badge, error) {
	return s.repo.ListBadges()
}

// GetBadgeByID finds a badge by ID
func (s *service) GetBadgeByID(id uint) (*Badge, error) {
	badge, err := s.repo.GetBadgeByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBadgeNotFound
		}
		return nil, err
	}
	return badge, nil
}

// AssignBadgeToTool assigns a badge to a tool
func (s *service) AssignBadgeToTool(toolID, badgeID uint) error {
	// Check badge exists
	_, err := s.repo.GetBadgeByID(badgeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBadgeNotFound
		}
		return err
	}

	// Check if already assigned
	exists, err := s.repo.BadgeExistsOnTool(toolID, badgeID)
	if err != nil {
		return err
	}
	if exists {
		return ErrBadgeAlreadyAssigned
	}

	return s.repo.AssignBadgeToTool(toolID, badgeID)
}

// RemoveBadgeFromTool removes a badge from a tool
func (s *service) RemoveBadgeFromTool(toolID, badgeID uint) error {
	// Check if badge is assigned
	exists, err := s.repo.BadgeExistsOnTool(toolID, badgeID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrBadgeNotAssigned
	}

	return s.repo.RemoveBadgeFromTool(toolID, badgeID)
}

// GetToolBadges returns all badges for a tool
func (s *service) GetToolBadges(toolID uint) ([]Badge, error) {
	return s.repo.GetToolBadges(toolID)
}

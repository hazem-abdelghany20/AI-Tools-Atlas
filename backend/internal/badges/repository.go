package badges

import (
	"gorm.io/gorm"
)

// Repository defines the interface for badge data operations
type Repository interface {
	ListBadges() ([]Badge, error)
	GetBadgeByID(id uint) (*Badge, error)
	Create(badge *Badge) error
	AssignBadgeToTool(toolID, badgeID uint) error
	RemoveBadgeFromTool(toolID, badgeID uint) error
	GetToolBadges(toolID uint) ([]Badge, error)
	BadgeExistsOnTool(toolID, badgeID uint) (bool, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new badge repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ListBadges returns all badges
func (r *repository) ListBadges() ([]Badge, error) {
	var badges []Badge
	err := r.db.Order("name ASC").Find(&badges).Error
	if err != nil {
		return nil, err
	}
	return badges, nil
}

// GetBadgeByID finds a badge by ID
func (r *repository) GetBadgeByID(id uint) (*Badge, error) {
	var badge Badge
	err := r.db.Where("id = ?", id).First(&badge).Error
	if err != nil {
		return nil, err
	}
	return &badge, nil
}

// Create inserts a new badge
func (r *repository) Create(badge *Badge) error {
	return r.db.Create(badge).Error
}

// AssignBadgeToTool assigns a badge to a tool
func (r *repository) AssignBadgeToTool(toolID, badgeID uint) error {
	return r.db.Exec(
		"INSERT INTO tool_badges (tool_id, badge_id, assigned_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT DO NOTHING",
		toolID, badgeID,
	).Error
}

// RemoveBadgeFromTool removes a badge from a tool
func (r *repository) RemoveBadgeFromTool(toolID, badgeID uint) error {
	return r.db.Exec("DELETE FROM tool_badges WHERE tool_id = ? AND badge_id = ?", toolID, badgeID).Error
}

// GetToolBadges returns all badges for a tool
func (r *repository) GetToolBadges(toolID uint) ([]Badge, error) {
	var badges []Badge
	err := r.db.Raw(`
		SELECT b.* FROM badges b
		JOIN tool_badges tb ON b.id = tb.badge_id
		WHERE tb.tool_id = ?
		ORDER BY b.name
	`, toolID).Scan(&badges).Error
	return badges, err
}

// BadgeExistsOnTool checks if a badge is assigned to a tool
func (r *repository) BadgeExistsOnTool(toolID, badgeID uint) (bool, error) {
	var count int64
	err := r.db.Raw("SELECT COUNT(*) FROM tool_badges WHERE tool_id = ? AND badge_id = ?", toolID, badgeID).Scan(&count).Error
	return count > 0, err
}

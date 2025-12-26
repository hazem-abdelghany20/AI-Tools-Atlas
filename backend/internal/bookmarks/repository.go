package bookmarks

import (
	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

// Repository defines the interface for bookmark data operations
type Repository interface {
	GetUserBookmarks(userID uint, sessionID string) ([]domain.Bookmark, error)
	AddBookmark(userID uint, sessionID string, toolID uint) (*domain.Bookmark, error)
	RemoveBookmark(userID uint, sessionID string, toolID uint) error
	IsBookmarked(userID uint, sessionID string, toolID uint) (bool, error)
	MigrateSessionBookmarks(userID uint, sessionID string) error
	UpdateToolBookmarkCount(toolID uint, delta int) error
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new bookmarks repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetUserBookmarks returns all bookmarks for a user or session
func (r *repository) GetUserBookmarks(userID uint, sessionID string) ([]domain.Bookmark, error) {
	var bookmarks []domain.Bookmark

	query := r.db.Model(&domain.Bookmark{}).
		Preload("Tool").
		Preload("Tool.PrimaryCategory").
		Preload("Tool.Tags").
		Preload("Tool.Badges")

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	} else {
		return []domain.Bookmark{}, nil
	}

	err := query.Order("created_at DESC").Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}

	return bookmarks, nil
}

// AddBookmark creates a new bookmark
func (r *repository) AddBookmark(userID uint, sessionID string, toolID uint) (*domain.Bookmark, error) {
	bookmark := &domain.Bookmark{
		ToolID: toolID,
	}

	if userID > 0 {
		bookmark.UserID = userID
	} else {
		bookmark.SessionID = sessionID
	}

	err := r.db.Create(bookmark).Error
	if err != nil {
		return nil, err
	}

	// Load the tool for the response
	err = r.db.
		Preload("Tool").
		Preload("Tool.PrimaryCategory").
		Preload("Tool.Tags").
		Preload("Tool.Badges").
		First(bookmark, bookmark.ID).Error
	if err != nil {
		return bookmark, nil // Return bookmark even if preload fails
	}

	return bookmark, nil
}

// RemoveBookmark deletes a bookmark
func (r *repository) RemoveBookmark(userID uint, sessionID string, toolID uint) error {
	query := r.db.Where("tool_id = ?", toolID)

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	} else {
		return gorm.ErrRecordNotFound
	}

	result := query.Delete(&domain.Bookmark{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// IsBookmarked checks if a tool is bookmarked by user or session
func (r *repository) IsBookmarked(userID uint, sessionID string, toolID uint) (bool, error) {
	var count int64

	query := r.db.Model(&domain.Bookmark{}).Where("tool_id = ?", toolID)

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	} else if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	} else {
		return false, nil
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// MigrateSessionBookmarks moves session bookmarks to a user account
func (r *repository) MigrateSessionBookmarks(userID uint, sessionID string) error {
	if userID == 0 || sessionID == "" {
		return nil
	}

	// Update all session bookmarks to belong to the user
	// Ignore duplicates (where user already has the bookmark)
	return r.db.Exec(`
		UPDATE bookmarks
		SET user_id = ?, session_id = NULL
		WHERE session_id = ?
		AND tool_id NOT IN (SELECT tool_id FROM bookmarks WHERE user_id = ?)
	`, userID, sessionID, userID).Error
}

// UpdateToolBookmarkCount updates the bookmark count for a tool
func (r *repository) UpdateToolBookmarkCount(toolID uint, delta int) error {
	return r.db.Model(&domain.Tool{}).
		Where("id = ?", toolID).
		UpdateColumn("bookmark_count", gorm.Expr("bookmark_count + ?", delta)).Error
}

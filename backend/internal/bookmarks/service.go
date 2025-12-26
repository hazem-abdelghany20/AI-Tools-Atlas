package bookmarks

import (
	"errors"

	"github.com/your-org/ai-tools-atlas-backend/internal/domain"
	"gorm.io/gorm"
)

// Error constants
var (
	ErrBookmarkNotFound = errors.New("bookmark not found")
	ErrAlreadyBookmarked = errors.New("tool already bookmarked")
	ErrInvalidRequest    = errors.New("invalid request: user_id or session_id required")
)

// BookmarkResponse represents a bookmark with tool info for API response
type BookmarkResponse struct {
	ID        uint        `json:"id"`
	ToolID    uint        `json:"tool_id"`
	Tool      domain.Tool `json:"tool"`
	CreatedAt string      `json:"created_at"`
}

// Service defines the interface for bookmark business logic
type Service interface {
	GetBookmarks(userID uint, sessionID string) ([]BookmarkResponse, error)
	AddBookmark(userID uint, sessionID string, toolID uint) (*BookmarkResponse, error)
	RemoveBookmark(userID uint, sessionID string, toolID uint) error
	IsBookmarked(userID uint, sessionID string, toolID uint) (bool, error)
	MigrateSessionBookmarks(userID uint, sessionID string) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new bookmarks service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// GetBookmarks returns all bookmarks for a user or session
func (s *service) GetBookmarks(userID uint, sessionID string) ([]BookmarkResponse, error) {
	if userID == 0 && sessionID == "" {
		return []BookmarkResponse{}, nil
	}

	bookmarks, err := s.repo.GetUserBookmarks(userID, sessionID)
	if err != nil {
		return nil, err
	}

	responses := make([]BookmarkResponse, len(bookmarks))
	for i, b := range bookmarks {
		responses[i] = s.toBookmarkResponse(b)
	}

	return responses, nil
}

// AddBookmark adds a tool to bookmarks
func (s *service) AddBookmark(userID uint, sessionID string, toolID uint) (*BookmarkResponse, error) {
	if userID == 0 && sessionID == "" {
		return nil, ErrInvalidRequest
	}

	// Check if already bookmarked
	isBookmarked, err := s.repo.IsBookmarked(userID, sessionID, toolID)
	if err != nil {
		return nil, err
	}
	if isBookmarked {
		return nil, ErrAlreadyBookmarked
	}

	// Add bookmark
	bookmark, err := s.repo.AddBookmark(userID, sessionID, toolID)
	if err != nil {
		return nil, err
	}

	// Update tool bookmark count - best effort, don't fail if this errors
	_ = s.repo.UpdateToolBookmarkCount(toolID, 1)

	resp := s.toBookmarkResponse(*bookmark)
	return &resp, nil
}

// RemoveBookmark removes a tool from bookmarks
func (s *service) RemoveBookmark(userID uint, sessionID string, toolID uint) error {
	if userID == 0 && sessionID == "" {
		return ErrInvalidRequest
	}

	err := s.repo.RemoveBookmark(userID, sessionID, toolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookmarkNotFound
		}
		return err
	}

	// Update tool bookmark count - best effort, don't fail if this errors
	_ = s.repo.UpdateToolBookmarkCount(toolID, -1)

	return nil
}

// IsBookmarked checks if a tool is bookmarked
func (s *service) IsBookmarked(userID uint, sessionID string, toolID uint) (bool, error) {
	return s.repo.IsBookmarked(userID, sessionID, toolID)
}

// MigrateSessionBookmarks moves session bookmarks to user account on login
func (s *service) MigrateSessionBookmarks(userID uint, sessionID string) error {
	return s.repo.MigrateSessionBookmarks(userID, sessionID)
}

// toBookmarkResponse converts a domain bookmark to API response
func (s *service) toBookmarkResponse(b domain.Bookmark) BookmarkResponse {
	return BookmarkResponse{
		ID:        b.ID,
		ToolID:    b.ToolID,
		Tool:      b.Tool,
		CreatedAt: b.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

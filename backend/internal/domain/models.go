package domain

import "time"

// Category represents an AI tool category
type Category struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Slug         string    `gorm:"uniqueIndex;not null" json:"slug"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description,omitempty"`
	IconURL      string    `gorm:"column:icon_url" json:"icon_url,omitempty"`
	DisplayOrder int       `gorm:"not null;default:0" json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Badge represents a badge/label that can be assigned to tools
type Badge struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Slug        string    `gorm:"uniqueIndex;not null" json:"slug"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	IconURL     string    `gorm:"column:icon_url" json:"icon_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Tag represents a tag that can be associated with tools
type Tag struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Slug      string    `gorm:"uniqueIndex;not null" json:"slug"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Media represents screenshots and videos for a tool
type Media struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ToolID       uint      `gorm:"not null" json:"tool_id"`
	Type         string    `gorm:"type:varchar(50);not null;check:type IN ('screenshot', 'video')" json:"type"`
	URL          string    `gorm:"type:varchar(500);not null" json:"url"`
	ThumbnailURL string    `gorm:"column:thumbnail_url;type:varchar(500)" json:"thumbnail_url,omitempty"`
	DisplayOrder int       `gorm:"not null;default:0" json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// Tool represents an AI tool
type Tool struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	Slug              string     `gorm:"uniqueIndex;not null" json:"slug"`
	Name              string     `gorm:"not null" json:"name"`
	LogoURL           string     `gorm:"column:logo_url" json:"logo_url,omitempty"`
	Tagline           string     `json:"tagline,omitempty"`
	Description       string     `gorm:"type:text" json:"description,omitempty"`
	BestFor           string     `gorm:"column:best_for;type:text" json:"best_for,omitempty"`
	PrimaryUseCases   string     `gorm:"type:text" json:"primary_use_cases,omitempty"`
	PricingSummary    string     `json:"pricing_summary,omitempty"`
	TargetRoles       string     `gorm:"type:text" json:"target_roles,omitempty"`
	Platforms         string     `gorm:"type:text" json:"platforms,omitempty"`
	HasFreeTier       bool       `gorm:"default:false" json:"has_free_tier"`
	OfficialURL       string     `gorm:"column:official_url" json:"official_url,omitempty"`
	PrimaryCategoryID uint       `json:"primary_category_id"`
	PrimaryCategory   Category   `gorm:"foreignKey:PrimaryCategoryID" json:"primary_category,omitempty"`
	AvgRatingOverall  float64    `gorm:"column:avg_rating_overall;default:0" json:"avg_rating_overall"`
	ReviewCount       int        `gorm:"default:0" json:"review_count"`
	BookmarkCount     int        `gorm:"default:0" json:"bookmark_count"`
	TrendingScore     float64    `gorm:"default:0" json:"trending_score"`
	Tags              []Tag      `gorm:"many2many:tool_tags" json:"tags,omitempty"`
	Media             []Media    `gorm:"foreignKey:ToolID" json:"media,omitempty"`
	Badges            []Badge    `gorm:"many2many:tool_badges" json:"badges,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	ArchivedAt        *time.Time `gorm:"index" json:"archived_at,omitempty"`
}

// ToolBadge is the join table for Tool-Badge many2many with extra field
type ToolBadge struct {
	ToolID     uint      `gorm:"primaryKey" json:"tool_id"`
	BadgeID    uint      `gorm:"primaryKey" json:"badge_id"`
	AssignedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"assigned_at"`
}

// ToolAlternative represents an alternative/similar tool relationship
type ToolAlternative struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ToolID            uint      `gorm:"not null" json:"tool_id"`
	AlternativeToolID uint      `gorm:"column:alternative_tool_id;not null" json:"alternative_tool_id"`
	RelationshipType  string    `gorm:"type:varchar(50);not null;check:relationship_type IN ('similar', 'alternative')" json:"relationship_type"`
	CreatedAt         time.Time `json:"created_at"`
}

// User represents a registered user
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"` // Never expose password hash in JSON
	DisplayName  string    `json:"display_name,omitempty"`
	Role         string    `gorm:"type:varchar(50);not null;check:role IN ('user', 'admin', 'moderator');default:'user'" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Review represents a user review of a tool
type Review struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	ToolID           uint       `gorm:"not null" json:"tool_id"`
	Tool             Tool       `gorm:"foreignKey:ToolID" json:"tool,omitempty"`
	UserID           uint       `gorm:"not null" json:"user_id"`
	User             User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RatingOverall    int        `gorm:"not null;check:rating_overall >= 1 AND rating_overall <= 5" json:"rating_overall"`
	RatingEaseOfUse  *int       `gorm:"check:rating_ease_of_use >= 1 AND rating_ease_of_use <= 5" json:"rating_ease_of_use,omitempty"`
	RatingValue      *int       `gorm:"check:rating_value >= 1 AND rating_value <= 5" json:"rating_value,omitempty"`
	RatingAccuracy   *int       `gorm:"check:rating_accuracy >= 1 AND rating_accuracy <= 5" json:"rating_accuracy,omitempty"`
	RatingSpeed      *int       `gorm:"check:rating_speed >= 1 AND rating_speed <= 5" json:"rating_speed,omitempty"`
	RatingSupport    *int       `gorm:"check:rating_support >= 1 AND rating_support <= 5" json:"rating_support,omitempty"`
	Pros             string     `gorm:"type:text" json:"pros,omitempty"`
	Cons             string     `gorm:"type:text" json:"cons,omitempty"`
	PrimaryUseCase   string     `json:"primary_use_case,omitempty"`
	ReviewerRole     string     `json:"reviewer_role,omitempty"`
	CompanySize      string     `json:"company_size,omitempty"`
	UsageContext     string     `gorm:"type:text" json:"usage_context,omitempty"`
	HelpfulCount     int        `gorm:"default:0" json:"helpful_count"`
	ModerationStatus string     `gorm:"type:varchar(50);not null;check:moderation_status IN ('pending', 'approved', 'rejected', 'hidden', 'removed');default:'pending'" json:"moderation_status"`
	ModeratedBy      *uint      `json:"moderated_by,omitempty"`
	ModeratedAt      *time.Time `json:"moderated_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Bookmark represents a user's saved tool
type Bookmark struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_tool" json:"user_id,omitempty"`
	SessionID string    `gorm:"uniqueIndex:idx_session_tool" json:"session_id,omitempty"`
	ToolID    uint      `gorm:"not null;uniqueIndex:idx_user_tool;uniqueIndex:idx_session_tool" json:"tool_id"`
	Tool      Tool      `gorm:"foreignKey:ToolID" json:"tool,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Report represents a content report (for tools or reviews)
type Report struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ReportableType string     `gorm:"type:varchar(50);not null;check:reportable_type IN ('tool', 'review')" json:"reportable_type"`
	ReportableID   uint       `gorm:"not null" json:"reportable_id"`
	ReporterUserID *uint      `json:"reporter_user_id,omitempty"`
	ReporterUser   *User      `gorm:"foreignKey:ReporterUserID" json:"reporter_user,omitempty"`
	Reason         string     `gorm:"type:varchar(50);not null;check:reason IN ('spam', 'abuse', 'misinformation', 'other')" json:"reason"`
	Comment        string     `gorm:"type:text" json:"comment,omitempty"`
	Status         string     `gorm:"type:varchar(50);not null;check:status IN ('pending', 'reviewed', 'dismissed');default:'pending'" json:"status"`
	ReviewedBy     *uint      `json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ModerationAction represents an audit log of moderation actions
type ModerationAction struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ReviewID    uint      `gorm:"not null" json:"review_id"`
	Review      Review    `gorm:"foreignKey:ReviewID" json:"-"`
	ModeratorID uint      `gorm:"not null" json:"moderator_id"`
	Moderator   User      `gorm:"foreignKey:ModeratorID" json:"moderator,omitempty"`
	ActionType  string    `gorm:"type:varchar(50);not null;check:action_type IN ('approve', 'hide', 'remove', 'restore')" json:"action_type"`
	Notes       string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName overrides for GORM
func (Category) TableName() string         { return "categories" }
func (Badge) TableName() string            { return "badges" }
func (Tag) TableName() string              { return "tags" }
func (Media) TableName() string            { return "media" }
func (Tool) TableName() string             { return "tools" }
func (ToolBadge) TableName() string        { return "tool_badges" }
func (ToolAlternative) TableName() string  { return "tool_alternatives" }
func (User) TableName() string             { return "users" }
func (Review) TableName() string           { return "reviews" }
func (Bookmark) TableName() string         { return "bookmarks" }
func (Report) TableName() string           { return "reports" }
func (ModerationAction) TableName() string { return "moderation_actions" }

### Story 1.3: GORM Models

**As a** developer,
**I want** GORM model structs defined for all entities,
**So that** I can perform type-safe database operations.

**Acceptance Criteria:**

**Given** I have the database schema from Story 1.2
**When** I create GORM model files
**Then** the following models are defined:

**In `internal/tools/model.go`:**
```go
type Tool struct {
    ID                 uint      `gorm:"primaryKey"`
    Slug               string    `gorm:"uniqueIndex;not null"`
    Name               string    `gorm:"not null"`
    LogoURL            string    `gorm:"column:logo_url"`
    Tagline            string
    Description        string    `gorm:"type:text"`
    BestFor            string    `gorm:"column:best_for;type:text"`
    PrimaryUseCases    string    `gorm:"type:text"`
    PricingSummary     string
    TargetRoles        string    `gorm:"type:text"`
    Platforms          string    `gorm:"type:text"`
    HasFreeTier        bool      `gorm:"default:false"`
    OfficialURL        string    `gorm:"column:official_url"`
    PrimaryCategoryID  uint
    PrimaryCategory    Category  `gorm:"foreignKey:PrimaryCategoryID"`
    AvgRatingOverall   float64   `gorm:"column:avg_rating_overall;default:0"`
    ReviewCount        int       `gorm:"default:0"`
    BookmarkCount      int       `gorm:"default:0"`
    TrendingScore      float64   `gorm:"default:0"`
    Tags               []Tag     `gorm:"many2many:tool_tags"`
    Media              []Media   `gorm:"foreignKey:ToolID"`
    Badges             []Badge   `gorm:"many2many:tool_badges"`
    CreatedAt          time.Time
    UpdatedAt          time.Time
    ArchivedAt         *time.Time `gorm:"index"`
}
```

**Similar models for:**
- `Category` in `internal/categories/model.go`
- `Tag` in `internal/tags/model.go` (future, can be in tools package initially)
- `Media` in `internal/tools/model.go` or separate media package
- `Badge` in `internal/badges/model.go`
- `User` in `internal/auth/model.go`
- `Review` in `internal/reviews/model.go`
- `Bookmark` in `internal/bookmarks/model.go`
- `ToolAlternative` in `internal/tools/model.go`

**And:**
- All models use proper GORM tags matching database schema
- JSON tags added for API serialization using snake_case: `json:"slug"`
- Relationships configured with appropriate foreign keys and constraints
- Auto-migration can run successfully: `db.AutoMigrate(&Tool{}, &Category{}, ...)` (for development)

**Technical Implementation:**

- PascalCase for struct names, camelCase for field names (Go convention)
- `gorm:` tags match database column names in snake_case
- `json:` tags use snake_case to match API response format
- Use pointers for nullable fields (`*time.Time`, `*string`)
- Omit empty fields in JSON with `json:"field,omitempty"` where appropriate

**Prerequisites:** Story 1.2 (Database Schema)

**Files Created:**
- `backend/internal/tools/model.go`
- `backend/internal/categories/model.go`
- `backend/internal/reviews/model.go`
- `backend/internal/bookmarks/model.go`
- `backend/internal/badges/model.go`
- `backend/internal/auth/model.go`

---

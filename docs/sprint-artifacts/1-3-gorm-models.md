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

### Tasks/Subtasks

- [x] Create directory structure for models
  - [x] Create `backend/internal/tools/` directory
  - [x] Create `backend/internal/categories/` directory
  - [x] Create `backend/internal/reviews/` directory
  - [x] Create `backend/internal/bookmarks/` directory
  - [x] Create `backend/internal/badges/` directory
  - [x] Create `backend/internal/auth/` directory
- [x] Define Tool model
  - [x] Create `internal/tools/model.go`
  - [x] Define Tool struct with all fields from schema
  - [x] Add GORM tags for database mapping
  - [x] Add JSON tags for API serialization (snake_case)
  - [x] Configure relationships (Category, Tags, Media, Badges)
- [x] Define Category model
  - [x] Create `internal/categories/model.go`
  - [x] Define Category struct with all fields
  - [x] Add appropriate tags
- [x] Define Tag model
  - [x] Add Tag struct to `internal/tools/model.go`
  - [x] Define with id, slug, name, timestamps
- [x] Define Media model
  - [x] Add Media struct to `internal/tools/model.go`
  - [x] Define with tool relationship and type enum
- [x] Define Badge model
  - [x] Create `internal/badges/model.go`
  - [x] Define Badge struct
- [x] Define User model
  - [x] Create `internal/auth/model.go`
  - [x] Define User struct with role enum
  - [x] Use pointer for nullable fields
- [x] Define Review model
  - [x] Create `internal/reviews/model.go`
  - [x] Define Review struct with all rating fields
  - [x] Configure moderation_status enum
- [x] Define Bookmark model
  - [x] Create `internal/bookmarks/model.go`
  - [x] Define with composite unique constraint
- [x] Define ToolAlternative model
  - [x] Add to `internal/tools/model.go`
  - [x] Configure relationship_type enum
- [x] Write tests for models
  - [x] Test GORM auto-migration works for all models
  - [x] Test relationships are properly configured
  - [x] Test JSON marshaling produces snake_case
- [x] Validate model definitions
  - [x] Run AutoMigrate successfully
  - [x] Verify generated SQL matches migration schema

---

### Dev Notes

**GORM Tag Conventions:**
- Use `gorm:"column:snake_case"` for non-standard column names
- Use `gorm:"uniqueIndex"` for unique constraints
- Use `gorm:"foreignKey:FieldName"` for explicit foreign keys
- Use `gorm:"many2many:join_table_name"` for many-to-many

**JSON Serialization:**
- All JSON tags must use snake_case to match API spec
- Use `json:"field,omitempty"` for optional fields

**Nullable Fields:**
- Use pointers (*string, *time.Time) for nullable database fields

---

### Dev Agent Record

#### Implementation Plan
1. Created all model packages: categories, tools, auth, reviews, bookmarks, badges
2. Defined all GORM models with proper tags and relationships
3. Created comprehensive test suite for auto-migration, relationships, and JSON serialization

#### Debug Log
- All models use snake_case JSON tags as required
- User model excludes password_hash from JSON with json:"-" tag
- Nullable fields use pointers (*int, *time.Time)
- Relationships configured with proper GORM tags

#### Completion Notes
✅ All 10 model types defined: Tool, Category, Tag, Media, Badge, ToolBadge, ToolAlternative, User, Review, Bookmark
✅ GORM tags match database schema from migrations
✅ JSON tags use snake_case for API compatibility
✅ Relationships configured: Tool->Category (belongs to), Tool->Tags (many2many), Tool->Media (has many), Tool->Badges (many2many)
✅ Comprehensive test suite validates:
- Auto-migration creates all tables
- Relationships work correctly (including Tool->Badges)
- JSON marshaling produces snake_case
- Unique constraints enforced
- Password hash not exposed in JSON

#### Code Review Fixes Applied (2025-12-26)
🔧 **CRITICAL:** Removed duplicate Badge model from tools/model.go - now uses badges.Badge only
🔧 **CRITICAL:** Fixed ErrorResponse calls in middleware.go (added missing statusCode parameter)
🔧 **HIGH:** Added ToolBadge join table model with AssignedAt field to capture tool_badges.assigned_at
🔧 **MEDIUM:** Added Tool->Badges many2many relationship test

---

### File List
- backend/internal/categories/model.go (created)
- backend/internal/tools/model.go (created, modified during code review)
- backend/internal/auth/model.go (created)
- backend/internal/reviews/model.go (created)
- backend/internal/bookmarks/model.go (created)
- backend/internal/badges/model.go (created)
- backend/internal/platform/db/models_test.go (created, enhanced during code review)
- backend/internal/platform/http/middleware.go (modified during code review - fixed ErrorResponse calls)

---

### Change Log
- 2025-12-26: Story 1-3 completed - All GORM models defined with tests
- 2025-12-26: Code review completed - Fixed duplicate Badge, added ToolBadge model, fixed middleware ErrorResponse calls, added Badges relationship test

---

### Status
Done

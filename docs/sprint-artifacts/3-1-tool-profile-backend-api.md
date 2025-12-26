### Story 3.1: Tool Profile Backend API

**As a** developer,
**I want** a complete tool profile endpoint with all related data,
**So that** users can view comprehensive tool information on the frontend.

**Acceptance Criteria:**

**Given** I need detailed tool profile data
**When** I implement the tool profile endpoint
**Then** the following endpoint works:

**Endpoint: `GET /api/v1/tools/:slug`**
- Accepts tool slug in URL path
- Returns complete tool object with all relationships loaded
- Response format: `{ data: { ...complete_tool_object } }`
- 404 error if tool slug not found or tool is archived

**Tool Object Includes:**
- **Core metadata:** id, slug, name, logo_url, tagline, description, best_for, primary_use_cases (parsed array or text), pricing_summary, target_roles (array or text), platforms (array or text), has_free_tier, official_url
- **Aggregated data:** avg_rating_overall, review_count, bookmark_count, trending_score
- **Relationships:**
  - `primary_category`: Full category object (id, slug, name, icon_url)
  - `tags`: Array of tag objects (id, slug, name)
  - `media`: Array of media objects (id, type, url, thumbnail_url, display_order) ordered by display_order
  - `badges`: Array of badge objects (id, slug, name, description, icon_url)
- **Timestamps:** created_at, updated_at

**And:**
- Repository function in `internal/tools/repository.go`:
  - `GetToolBySlug(slug string) (*Tool, error)` with all preloads
- Service function validates tool exists and is not archived
- Handler uses standardized success response

**Technical Implementation:**

- GORM query with comprehensive preloading:
  ```go
  db.Where("slug = ? AND archived_at IS NULL", slug).
     Preload("PrimaryCategory").
     Preload("Tags", func(db *gorm.DB) *gorm.DB {
         return db.Order("tags.name ASC")
     }).
     Preload("Media", func(db *gorm.DB) *gorm.DB {
         return db.Order("media.display_order ASC")
     }).
     Preload("Badges").
     First(&tool)
  ```
- Parse text fields into arrays if stored as comma-separated or JSON
- Return 404 with error envelope if not found

**Prerequisites:** Epic 1 (Foundation, models), Epic 2 (Tools repository base)

**Files Modified:**
- `backend/internal/tools/repository.go` (add GetToolBySlug)
- `backend/internal/tools/service.go` (add GetTool)
- `backend/internal/tools/handler.go` (add GET /:slug handler)

---

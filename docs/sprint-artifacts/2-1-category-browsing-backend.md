### Story 2.1: Category Browsing Backend

**As a** developer,
**I want** category listing and category-filtered tools endpoints implemented,
**So that** users can browse tools by category.

**Acceptance Criteria:**

**Given** I need category browsing functionality
**When** I implement the categories API
**Then** the following endpoints work:

**Endpoint: `GET /api/v1/categories`**
- Returns list of all categories with fields: id, slug, name, description, icon_url, display_order
- Response format: `{ data: [{ id, slug, name, description, icon_url, display_order }] }`
- Ordered by display_order ASC
- Only active (non-archived) categories returned

**Endpoint: `GET /api/v1/categories/:slug/tools`**
- Accepts category slug in URL path
- Returns paginated list of tools in that category
- Response format: `{ data: [...tools], meta: { page, page_size, total } }`
- Each tool includes: id, slug, name, logo_url, tagline, best_for, avg_rating_overall, review_count, bookmark_count, pricing_summary, has_free_tier, tags (array)
- Default pagination: page_size=20
- Query params: `?page=1&page_size=20`
- 404 error if category slug not found

**And:**
- Repository layer in `internal/categories/repository.go` with functions:
  - `ListCategories() ([]Category, error)`
  - `GetCategoryBySlug(slug string) (*Category, error)`
  - `ListToolsByCategory(categoryID uint, page, pageSize int) ([]Tool, int, error)`
- Service layer in `internal/categories/service.go` handling business logic
- Handler in `internal/categories/handler.go` using standardized responses
- Proper error handling with appropriate HTTP status codes

**Technical Implementation:**

- GORM queries with preloading for tags: `.Preload("Tags").Preload("PrimaryCategory")`
- Pagination: `.Limit(pageSize).Offset((page - 1) * pageSize)`
- Filter out archived tools: `WHERE archived_at IS NULL`
- Count query for total: `db.Model(&Tool{}).Where("primary_category_id = ?", categoryID).Count(&total)`

**Prerequisites:** Epic 1 (Foundation, GORM models)

**Files Created:**
- `backend/internal/categories/repository.go`
- `backend/internal/categories/service.go`
- `backend/internal/categories/handler.go`

---

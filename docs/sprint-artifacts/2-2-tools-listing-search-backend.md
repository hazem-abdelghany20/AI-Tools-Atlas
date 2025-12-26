### Story 2.2: Tools Listing & Search Backend

**As a** developer,
**I want** tools listing and free-text search endpoints implemented,
**So that** users can search and browse all tools.

**Acceptance Criteria:**

**Given** I need tools listing and search
**When** I implement the tools API
**Then** the following endpoints work:

**Endpoint: `GET /api/v1/tools`**
- Returns paginated list of all tools
- Response format: `{ data: [...tools], meta: { page, page_size, total } }`
- Each tool includes: id, slug, name, logo_url, tagline, best_for, avg_rating_overall, review_count, bookmark_count, pricing_summary, has_free_tier, primary_category (nested object), tags (array), badges (array)
- Query params:
  - `?page=1&page_size=20` (pagination)
  - `?category=slug` (filter by category)
  - `?price=free|freemium|paid` (filter by pricing)
  - `?min_rating=4` (filter by minimum rating)
  - `?platform=web|mobile|api` (filter by platform - search in platforms text field)
  - `?sort=top_rated|most_bookmarked|trending|newest` (sort order)
- Default sort: top_rated
- Only non-archived tools returned

**Endpoint: `GET /api/v1/search/tools`**
- Free-text search across tool name, tagline, description, best_for, primary_use_cases
- Query params: `?q=search+query` plus all filters from GET /tools
- Returns same format as GET /tools with matching tools
- Search is case-insensitive
- Empty query returns all tools (same as GET /tools)

**Sorting Logic:**
- `top_rated`: ORDER BY avg_rating_overall DESC, review_count DESC
- `most_bookmarked`: ORDER BY bookmark_count DESC
- `trending`: ORDER BY trending_score DESC (trending_score calculated separately)
- `newest`: ORDER BY created_at DESC

**And:**
- Repository in `internal/tools/repository.go`:
  - `ListTools(filters ToolFilters, page, pageSize int) ([]Tool, int, error)`
  - `SearchTools(query string, filters ToolFilters, page, pageSize int) ([]Tool, int, error)`
- Service in `internal/tools/service.go`
- Handler in `internal/tools/handler.go`
- ToolFilters struct: `{ Category, Price, MinRating, Platform, Sort }`

**Technical Implementation:**

- Full-text search using PostgreSQL ILIKE or tsvector for better performance:
  ```sql
  WHERE name ILIKE '%query%' OR tagline ILIKE '%query%' OR description ILIKE '%query%' OR best_for ILIKE '%query%'
  ```
- Preload relationships: `.Preload("PrimaryCategory").Preload("Tags").Preload("Badges")`
- Dynamic query building based on filters:
  ```go
  query := db.Model(&Tool{}).Where("archived_at IS NULL")
  if filters.Category != "" {
      query = query.Joins("JOIN categories ON categories.id = tools.primary_category_id").Where("categories.slug = ?", filters.Category)
  }
  if filters.MinRating > 0 {
      query = query.Where("avg_rating_overall >= ?", filters.MinRating)
  }
  // ... etc
  ```

**Prerequisites:** Epic 1, Story 2.1 (Categories)

**Files Created:**
- `backend/internal/tools/repository.go`
- `backend/internal/tools/service.go`
- `backend/internal/tools/handler.go`

---

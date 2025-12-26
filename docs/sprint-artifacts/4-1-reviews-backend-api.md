### Story 4.1: Reviews Backend API

**As a** developer,
**I want** endpoints for fetching and submitting reviews,
**So that** users can view and contribute structured reviews.

**Acceptance Criteria:**

**Endpoint: `GET /api/v1/tools/:slug/reviews`**
- Returns paginated list of reviews for a tool
- Response: `{ data: [...reviews], meta: { page, page_size, total } }`
- Each review includes: id, rating_overall, rating_ease_of_use, rating_value, rating_accuracy, rating_speed, rating_support, pros, cons, primary_use_case, reviewer_role, company_size, usage_context, helpful_count, created_at, user (id, display_name)
- Query params: `?page=1&page_size=10&sort=newest|most_helpful|highest|lowest`
- Default sort: newest
- Only approved/non-hidden reviews returned

**Endpoint: `POST /api/v1/tools/:slug/reviews`** (Authenticated)
- Accepts review submission from logged-in users
- Request body: rating_overall (required), rating dimensions (optional), pros (required), cons (required), primary_use_case, reviewer_role, company_size, usage_context
- Validation: rating 1-5, pros/cons max 500 chars each, required fields present
- Creates review with moderation_status: auto-approved or pending based on config
- Updates tool.avg_rating_overall and tool.review_count
- Returns created review object
- 401 if not authenticated, 422 if validation fails

**And:**
- Repository: `ListReviewsByTool`, `CreateReview`, `UpdateToolRatingAggregates`
- Service: validates input, checks user hasn't already reviewed this tool
- Handler: uses auth middleware, standardized responses

**Technical Implementation:**
- GORM preload user: `.Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id", "display_name") })`
- Rating aggregation: recalculate avg_rating_overall from all approved reviews
- Rate limiting: prevent spam, 1 review per user per tool

**Prerequisites:** Epic 1, Epic 3 (Tool profile)

**Files Modified:**
- `backend/internal/reviews/repository.go`
- `backend/internal/reviews/service.go`
- `backend/internal/reviews/handler.go`

---

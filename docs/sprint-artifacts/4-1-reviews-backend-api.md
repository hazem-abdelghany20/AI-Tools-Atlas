### Story 4.1: Reviews Backend API

**As a** developer,
**I want** endpoints for fetching and submitting reviews,
**So that** users can view and contribute structured reviews.

**Status:** done

---

## Acceptance Criteria:

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

---

## Tasks/Subtasks:

- [x] **Task 1: Create Reviews Repository**
  - [x] 1.1 Define Repository interface with ListReviewsByTool, CreateReview, HasUserReviewed, UpdateToolRatingAggregates
  - [x] 1.2 Implement ListReviewsByTool with pagination, sorting, and user preload
  - [x] 1.3 Implement CreateReview with proper associations
  - [x] 1.4 Implement HasUserReviewed to check duplicate reviews
  - [x] 1.5 Implement UpdateToolRatingAggregates to recalculate avg_rating and review_count

- [x] **Task 2: Create Reviews Service**
  - [x] 2.1 Define Service interface with ListReviews, CreateReview
  - [x] 2.2 Implement input validation (rating range, text length limits)
  - [x] 2.3 Implement duplicate review check (one review per user per tool)
  - [x] 2.4 Implement CreateReview with moderation status handling
  - [x] 2.5 Wire up rating aggregation after review creation

- [x] **Task 3: Create Reviews Handler**
  - [x] 3.1 Define Handler struct with service dependency
  - [x] 3.2 Implement GET /api/v1/tools/:slug/reviews endpoint
  - [x] 3.3 Implement POST /api/v1/tools/:slug/reviews endpoint with auth middleware
  - [x] 3.4 Parse and validate query parameters for sorting and pagination
  - [x] 3.5 Return standardized responses (success, error, validation errors)

- [x] **Task 4: Register Routes and Integration**
  - [x] 4.1 Update router.go to include reviews routes
  - [x] 4.2 Apply AuthRequired middleware to POST endpoint
  - [x] 4.3 Ensure proper error handling and status codes

- [x] **Task 5: Write Unit Tests**
  - [x] 5.1 Create service_test.go with mock repository
  - [x] 5.2 Test ListReviews with various sort options
  - [x] 5.3 Test CreateReview validation logic
  - [x] 5.4 Test duplicate review prevention
  - [x] 5.5 Create handler_test.go with mock service

---

## Dev Notes:

**Technical Implementation:**
- GORM preload user: `.Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id", "display_name") })`
- Rating aggregation: recalculate avg_rating_overall from all approved reviews
- Rate limiting: prevent spam, 1 review per user per tool

**Prerequisites:** Epic 1, Epic 3 (Tool profile)

**Files to Modify:**
- `backend/internal/reviews/repository.go`
- `backend/internal/reviews/service.go`
- `backend/internal/reviews/handler.go`
- `backend/internal/reviews/service_test.go`
- `backend/internal/reviews/handler_test.go`
- `backend/internal/platform/http/router.go`

---

## Dev Agent Record:

### Implementation Plan:
1. Create repository with GORM operations for reviews
2. Create service with business logic and validation
3. Create handler with HTTP endpoints
4. Register routes with auth middleware
5. Write comprehensive tests

### Debug Log:
- Started implementation: 2025-12-26

### Completion Notes:
- Implemented reviews repository with ListReviewsByTool, CreateReview, HasUserReviewed, UpdateToolRatingAggregates
- Implemented reviews service with validation for ratings (1-5), text length (500 chars), duplicate prevention
- Implemented reviews handler with GET and POST endpoints
- Added ReviewResponse DTO to include user info (id, display_name)
- Registered routes in router.go with auth middleware for POST
- Created comprehensive unit tests for service layer
- All tests pass

---

## File List:
- `backend/internal/reviews/repository.go` (created)
- `backend/internal/reviews/service.go` (created)
- `backend/internal/reviews/handler.go` (created)
- `backend/internal/reviews/service_test.go` (created)
- `backend/internal/reviews/handler_test.go` (created)
- `backend/internal/platform/http/router.go` (modified)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation started |
| 2025-12-26 | Created reviews repository, service, handler, tests |
| 2025-12-26 | Registered routes with auth middleware |
| 2025-12-26 | Code Review: Fixed unsafe type assertion in handler, UTF-8 char counting, empty CreatedAt, added handler_test.go |

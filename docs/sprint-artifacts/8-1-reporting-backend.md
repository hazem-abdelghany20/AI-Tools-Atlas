### Story 8.1: Reporting Backend

**As a** developer,
**I want** endpoints for reporting content,
**So that** users can flag problematic tools or reviews.

**Acceptance Criteria:**

**Endpoints:**

**`POST /api/v1/tools/:slug/report`**
- Request: `{ reason: "spam|abuse|misinformation|other", comment?: string }`
- Creates report record
- Returns 201 Created

**`POST /api/v1/reviews/:id/report`**
- Same structure as tool reports
- Creates report record
- Returns 201 Created

**Reports Table:**
- Fields: id, reportable_type (tool/review), reportable_id, reporter_user_id (nullable for anonymous), reason, comment, status (pending/reviewed/dismissed), created_at

**And:**
- Repository: CreateReport
- Service: validates reason enum, limits reports per user
- Both authenticated and anonymous users can report

**Prerequisites:** Epic 1

**New Files:**
- `backend/internal/moderation/repository.go`
- `backend/internal/moderation/service.go`
- `backend/internal/moderation/handler.go`

---

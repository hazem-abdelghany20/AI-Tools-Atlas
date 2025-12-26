### Story 8.3: Moderation Queue Backend ✅ DONE

**As a** developer,
**I want** moderation queue endpoints,
**So that** moderators can review reported content.

**Acceptance Criteria:**

**Endpoints (require moderator or admin role):**

**`GET /api/v1/moderation/queue`**
- Returns paginated list of reports
- Query params: `?type=tool|review&status=pending|reviewed|dismissed&page=1`
- Response: `{ data: [...reports], meta: { page, page_size, total } }`
- Each report includes: reportable object (tool or review with full details), reporter info, reason, comment, created_at

**`PATCH /api/v1/moderation/reviews/:id/approve`**
- Sets review.moderation_status = "approved"
- Logs moderation action
- Returns updated review

**`PATCH /api/v1/moderation/reviews/:id/hide`**
- Sets review.moderation_status = "hidden"
- Review no longer appears in public endpoints
- Logs moderation action

**`PATCH /api/v1/moderation/reviews/:id/remove`**
- Sets review.moderation_status = "removed"
- Permanently hides review
- Updates tool rating aggregates
- Logs moderation action

**`GET /api/v1/moderation/history/:review_id`**
- Returns audit log of moderation actions for a review
- Response: `{ data: [...actions] }` with actor, action_type, timestamp, notes

**And:**
- Moderation actions table: tracks who did what when
- Service layer validates moderator permissions
- Aggregates: recalculate tool ratings when reviews hidden/removed

**Prerequisites:** Story 8.1

**Files Modified:**
- `backend/internal/moderation/repository.go`
- `backend/internal/moderation/service.go`
- `backend/internal/moderation/handler.go`
- `backend/internal/reviews/service.go` (moderation status updates)

---

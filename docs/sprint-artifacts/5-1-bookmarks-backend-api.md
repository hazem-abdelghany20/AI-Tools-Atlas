### Story 5.1: Bookmarks Backend API

**As a** developer,
**I want** endpoints for managing user bookmarks,
**So that** users can save tools to their shortlist.

**Status:** done

---

## Acceptance Criteria:

**Endpoint: `GET /api/v1/me/bookmarks`** (Authenticated or session-based)
- Returns list of bookmarked tools for current user
- Response: `{ data: [...tools] }` (full tool objects)
- Ordered by bookmark created_at DESC

**Endpoint: `POST /api/v1/me/bookmarks`** (Authenticated or session-based)
- Request: `{ tool_id: number }`
- Creates bookmark for user/session
- Increments tool.bookmark_count
- Returns 201 with created bookmark
- Returns 409 if already bookmarked

**Endpoint: `DELETE /api/v1/me/bookmarks/:tool_id`** (Authenticated or session-based)
- Deletes bookmark for user/session
- Decrements tool.bookmark_count
- Returns 204 No Content

---

## Tasks/Subtasks:

- [x] **Task 1: Create Bookmarks Repository**
  - [x] 1.1 Define Repository interface
  - [x] 1.2 Implement GetUserBookmarks with tool preload
  - [x] 1.3 Implement AddBookmark
  - [x] 1.4 Implement RemoveBookmark
  - [x] 1.5 Implement IsBookmarked
  - [x] 1.6 Implement MigrateSessionBookmarks
  - [x] 1.7 Implement UpdateToolBookmarkCount

- [x] **Task 2: Create Bookmarks Service**
  - [x] 2.1 Define Service interface
  - [x] 2.2 Implement GetBookmarks
  - [x] 2.3 Implement AddBookmark with duplicate check
  - [x] 2.4 Implement RemoveBookmark

- [x] **Task 3: Create Bookmarks Handler**
  - [x] 3.1 Implement GET endpoint
  - [x] 3.2 Implement POST endpoint
  - [x] 3.3 Implement DELETE endpoint
  - [x] 3.4 Handle session ID generation

- [x] **Task 4: Update Infrastructure**
  - [x] 4.1 Update domain Bookmark model with SessionID
  - [x] 4.2 Add OptionalAuth middleware
  - [x] 4.3 Register routes in router.go

---

## File List:
- `backend/internal/bookmarks/repository.go` (created)
- `backend/internal/bookmarks/service.go` (created)
- `backend/internal/bookmarks/handler.go` (created)
- `backend/internal/domain/models.go` (modified)
- `backend/internal/platform/http/middleware.go` (modified)
- `backend/internal/platform/http/router.go` (modified)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation completed |
| 2025-12-26 | Code Review: Fixed cookie security (HttpOnly, Secure), improved error handling |

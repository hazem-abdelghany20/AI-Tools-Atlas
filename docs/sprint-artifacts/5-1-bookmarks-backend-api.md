### Story 5.1: Bookmarks Backend API

**As a** developer,
**I want** endpoints for managing user bookmarks,
**So that** users can save tools to their shortlist.

**Acceptance Criteria:**

**Endpoint: `GET /api/v1/me/bookmarks`** (Authenticated or session-based)
- Returns list of bookmarked tools for current user
- Response: `{ data: [...tools] }` (full tool objects, same as search results)
- Ordered by bookmark created_at DESC (most recent first)
- Supports anonymous users via session_id stored in bookmarks table

**Endpoint: `POST /api/v1/me/bookmarks`** (Authenticated or session-based)
- Request: `{ tool_id: number }`
- Creates bookmark for user/session
- Increments tool.bookmark_count
- Returns 201 with created bookmark
- Returns 409 if already bookmarked (idempotent)

**Endpoint: `DELETE /api/v1/me/bookmarks/:tool_id`** (Authenticated or session-based)
- Deletes bookmark for user/session
- Decrements tool.bookmark_count
- Returns 204 No Content
- Returns 404 if bookmark doesn't exist

**And:**
- Support anonymous sessions: generate session_id, store in cookie
- On user login: migrate session bookmarks to user account
- Repository: `GetUserBookmarks`, `AddBookmark`, `RemoveBookmark`
- Unique constraint on (user_id, tool_id) and (session_id, tool_id)

**Technical Implementation:**
- Session middleware: generate UUID session_id if not present
- Migration on login: UPDATE bookmarks SET user_id = ?, session_id = NULL WHERE session_id = ?
- Update tool counts: db.Model(&Tool{}).Where("id = ?", toolID).UpdateColumn("bookmark_count", gorm.Expr("bookmark_count + ?", delta))

**Prerequisites:** Epic 1 (Auth), Epic 2 (Tools)

**Files Created:**
- `backend/internal/bookmarks/repository.go`
- `backend/internal/bookmarks/service.go`
- `backend/internal/bookmarks/handler.go`

---

### Story 7.1: Admin Tools Management Backend

**Status:** Done ✅

**As a** developer,
**I want** admin endpoints for managing tools,
**So that** admins can maintain the catalog.

**Acceptance Criteria:**

**Endpoints (all require admin role):**

**`GET /api/v1/admin/tools`**
- Returns all tools including archived
- Query params: `?search=...&archived=true|false&page=1`
- Response: `{ data: [...tools], meta: { page, page_size, total } }`

**`POST /api/v1/admin/tools`**
- Create new tool
- Request: all tool fields (name, slug, description, best_for, category_id, etc.)
- Validation: required fields, unique slug
- Returns created tool

**`PATCH /api/v1/admin/tools/:id`**
- Update tool fields
- Request: partial tool object
- Validates changes
- Returns updated tool

**`DELETE /api/v1/admin/tools/:id`**
- Soft delete: sets archived_at timestamp
- Tool no longer appears in public endpoints
- Returns 204 No Content

**And:**
- AdminRequired middleware on all routes
- Service layer validates admin permissions
- Audit logging: track who made changes (future enhancement)

**Prerequisites:** Epic 1 (Auth with admin role)

**Files Modified:**
- `backend/internal/tools/handler.go` (add admin routes)
- `backend/internal/tools/service.go` (add admin methods)
- `backend/internal/platform/http/router.go` (register admin routes)

---

## Tasks

### Task 1: Repository Admin Methods ✅
- [x] Add `ListToolsAdmin` for listing tools with archived option
- [x] Add `GetToolByIDAdmin` for fetching tool including archived
- [x] Add `Create`, `Update`, `Archive` methods
- [x] Add `SlugExists` for validation
- [x] Update mock implementations in tests

### Task 2: Service Admin Methods ✅
- [x] Add `CreateToolInput` and `UpdateToolInput` DTOs
- [x] Add error constants for validation
- [x] Add `ListToolsAdmin`, `GetToolByIDAdmin`, `CreateTool`, `UpdateTool`, `ArchiveTool` methods
- [x] Implement validation logic

### Task 3: Handler Admin Endpoints ✅
- [x] Add `RegisterAdminRoutes` method
- [x] Implement `AdminListTools`, `AdminGetTool`, `AdminCreateTool`, `AdminUpdateTool`, `AdminArchiveTool`
- [x] Add proper error responses with status codes

### Task 4: Router Integration ✅
- [x] Add admin routes group with auth + admin middleware
- [x] Register tool admin routes

---

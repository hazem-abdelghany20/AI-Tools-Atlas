### Story 7.3: Admin Categories & Tags Management

**As an** admin,
**I want** to manage categories and tags,
**So that** taxonomy remains coherent.

**Acceptance Criteria:**

**Endpoints:**
- `GET/POST/PATCH/DELETE /api/v1/admin/categories`
- `GET/POST/PATCH/DELETE /api/v1/admin/tags`
- All require admin role

**Admin Categories Page:**
- List of categories with: Name, Slug, Tool Count, Display Order, Actions
- Create/Edit category: name, slug, description, icon_url, display_order
- Cannot delete category if tools exist (show warning, require reassignment)
- Drag-and-drop reordering by display_order

**Admin Tags Page:**
- List of tags with: Name, Slug, Tool Count, Actions
- Create/Edit tag: name, slug
- Delete tag: removes from all tools (confirmation required)

**Technical Implementation:**
- Similar pattern to tools management
- Category/tag CRUD operations
- Prevent deletion with referential integrity checks

**Prerequisites:** Story 7.1

**New Files:**
- `frontend/src/views/admin/AdminCategoriesView.vue`
- `frontend/src/views/admin/AdminTagsView.vue`

**Files Modified:**
- `backend/internal/categories/handler.go`
- Backend tag handlers (new files if not existing)

---

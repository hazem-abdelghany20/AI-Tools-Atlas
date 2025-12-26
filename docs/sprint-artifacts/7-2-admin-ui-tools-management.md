### Story 7.2: Admin UI - Tools Management

**Status:** Done ✅

**As an** admin,
**I want** a UI to manage tools,
**So that** I can keep the catalog accurate and complete.

**Acceptance Criteria:**

**Admin Tools Page (`/admin/tools`):**
- Table of all tools with columns: Logo, Name, Category, Rating, Reviews, Bookmarks, Status (Active/Archived), Actions
- Search bar: filter by name/slug
- Filter: All / Active / Archived
- "Add Tool" button opens tool form modal
- Actions per tool: Edit, Archive/Unarchive, View on site
- Pagination

**Tool Form Modal:**
- All fields editable: name, slug, tagline, description, best_for, primary_use_cases, pricing_summary, target_roles, platforms, has_free_tier, official_url, category, tags, media URLs
- Validation: required fields highlighted
- Preview: shows how tool card will look
- Submit: saves tool, closes modal, refreshes list
- Cancel: closes without saving

**Visual Design:**
- Admin UI uses same dark theme but with admin-specific nav
- Table is sortable and filterable
- Clear visual difference between active/archived tools

**Technical Implementation:**
- AdminToolsView.vue with data table
- ToolFormModal.vue component for create/edit
- Admin routes protected: check user role in router beforeEach
- API calls to admin endpoints

**Prerequisites:** Story 7.1 (Admin tools API)

**New Files:**
- `frontend/src/views/admin/AdminToolsView.vue`
- `frontend/src/components/admin/ToolFormModal.vue`
- `frontend/src/router/admin.ts` (admin routes)

---

## Tasks

### Task 1: Admin Tools View ✅
- [x] Create `AdminToolsView.vue` with data table
- [x] Implement search bar with debounce
- [x] Add status filter (All/Active/Archived)
- [x] Add pagination controls
- [x] Add admin role check on mount

### Task 2: Tool Form Modal ✅
- [x] Create `ToolFormModal.vue` component
- [x] Implement form for all tool fields
- [x] Add validation for required fields
- [x] Handle create and update modes
- [x] Add error display

### Task 3: Router Integration ✅
- [x] Add admin/tools route with meta flags
- [x] Add Admin Panel link in header for admin users
- [x] Add profile link in header dropdown

### Task 4: Types Update ✅
- [x] Add `archived_at` field to Tool type
- [x] Add `CreateToolInput` and `UpdateToolInput` types

---

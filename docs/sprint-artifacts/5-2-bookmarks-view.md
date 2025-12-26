### Story 5.2: Bookmarks View

**As a** user,
**I want** to view and manage my bookmarked tools,
**So that** I can access my shortlist and make decisions.

**Status:** done

---

## Acceptance Criteria:

**BookmarksView (`/bookmarks`):**
- Page title: "المفضلة" (My Bookmarks)
- Grid of bookmarked tools using ToolCard component
- Each card shows bookmark button (filled/active state)
- "Remove from bookmarks" action on each card
- "Compare selected" button: allows multi-select and send to comparison
- Empty state: "لم تحفظ أي أدوات بعد. ابدأ بالتصفح!" with link to home
- Anonymous users: bookmarks stored in localStorage, prompt to sign in to sync across devices

**Multi-Select for Comparison:**
- Checkbox on each card for selection
- "مقارنة المحددة" (Compare Selected) button appears when 2+ tools selected
- Click button navigates to `/compare?tools=slug1,slug2,slug3`
- Max 4 tools can be selected

**Technical Implementation:**
- Fetch bookmarks: `apiClient.get('/me/bookmarks')`
- useBookmarksStore to manage state
- Multi-select: local state tracking selected tool IDs
- Remove bookmark: calls store.removeBookmark(toolId), updates UI optimistically

**Prerequisites:** Story 5.1 (Bookmarks API), Story 2.4 (ToolCard)

---

## Tasks/Subtasks:

- [x] **Task 1: Create BookmarksView Component**
  - [x] 1.1 Create view with RTL layout
  - [x] 1.2 Implement header with title and description
  - [x] 1.3 Add loading state with skeleton cards
  - [x] 1.4 Add empty state with icon and CTA
  - [x] 1.5 Implement tools grid with ToolCard

- [x] **Task 2: Integrate Bookmarks Store**
  - [x] 2.1 Connect to useBookmarksStore
  - [x] 2.2 Fetch bookmarks on mount
  - [x] 2.3 Display bookmarked tools in grid

- [x] **Task 3: Multi-Select for Comparison**
  - [x] 3.1 Add local compareTools state
  - [x] 3.2 Track selection via ToolCard props
  - [x] 3.3 Handle add/remove from compare
  - [x] 3.4 Show CompareBar when tools selected

- [x] **Task 4: Route Setup**
  - [x] 4.1 Add /bookmarks route in router
  - [x] 4.2 Lazy load BookmarksView component

---

## File List:
- `frontend/src/views/BookmarksView.vue` (created)
- `frontend/src/router/index.ts` (modified)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation completed |
| 2025-12-26 | Code Review: Fixed compare route param (tools=slugs), added optimistic UI update |

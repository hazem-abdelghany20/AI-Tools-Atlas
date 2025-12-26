### Story 5.2: Bookmarks View

**As a** user,
**I want** to view and manage my bookmarked tools,
**So that** I can access my shortlist and make decisions.

**Acceptance Criteria:**

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

**Files Created:**
- `frontend/src/views/BookmarksView.vue`

---

### Story 5.3: Add to Compare from Results

**As a** user,
**I want** to add tools to comparison from search results and profiles,
**So that** I can build up a comparison set.

**Status:** done

---

## Acceptance Criteria:

**Compare Selection:**
- ToolCard component has "مقارنة" (Compare) checkbox/button
- Clicking adds tool to comparison set (max 4 tools)
- Visual indicator shows tool is in comparison set
- Floating comparison bar appears at bottom when 1+ tools selected:
  - Shows selected tool logos/names
  - "مقارنة" (Compare) button to go to comparison view
  - "مسح" (Clear) button to clear selection
  - Tool count: "2 من 4 أدوات"
- Clicking "Compare" button navigates to `/compare?tools=slug1,slug2`

**Persistence:**
- Comparison selection stored in Pinia store
- Persists across page navigation within session
- Option to save comparison set to bookmarks (future)

**Technical Implementation:**
- useComparisonStore: `{ selectedTools: [], addTool(), removeTool(), clear() }`
- ComparisonBar.vue component fixed at bottom
- ToolCard emits @add-to-compare event, parent calls store.addTool()
- Max 4 validation: show toast if limit reached

**Prerequisites:** Story 2.4 (ToolCard)

---

## Tasks/Subtasks:

- [x] **Task 1: Create Comparison Store**
  - [x] 1.1 Define store with selectedTools state
  - [x] 1.2 Implement addTool with max 4 limit
  - [x] 1.3 Implement removeTool
  - [x] 1.4 Implement clear
  - [x] 1.5 Implement isInComparison helper
  - [x] 1.6 Implement getCompareUrl helper
  - [x] 1.7 Add persist option for localStorage

- [x] **Task 2: Create CompareBar Component**
  - [x] 2.1 Create fixed bottom bar with RTL
  - [x] 2.2 Display selected tool pills with logos
  - [x] 2.3 Add remove button on each pill
  - [x] 2.4 Add Clear All button
  - [x] 2.5 Add Compare button (disabled if < 2)
  - [x] 2.6 Emit events for actions

- [x] **Task 3: Integrate in Views**
  - [x] 3.1 Update SearchResultsView with store and CompareBar
  - [x] 3.2 Update HomeView with store and CompareBar
  - [x] 3.3 Update BookmarksView (already has local compare)

- [x] **Task 4: ToolCard Integration**
  - [x] 4.1 Verify ToolCard has compare button
  - [x] 4.2 Verify isInCompare prop works
  - [x] 4.3 Verify events are emitted correctly

---

## File List:
- `frontend/src/stores/comparison.ts` (created)
- `frontend/src/components/common/CompareBar.vue` (created)
- `frontend/src/views/SearchResultsView.vue` (modified)
- `frontend/src/views/HomeView.vue` (modified)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Story implementation completed |
| 2025-12-26 | Code Review: Fixed getCompareUrl to use slugs per AC |

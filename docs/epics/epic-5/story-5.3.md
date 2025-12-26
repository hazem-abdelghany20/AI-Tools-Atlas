### Story 5.3: Add to Compare from Results

**As a** user,
**I want** to add tools to comparison from search results and profiles,
**So that** I can build up a comparison set.

**Acceptance Criteria:**

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

**New Files:**
- `frontend/src/stores/comparison.ts`
- `frontend/src/components/compare/ComparisonBar.vue`

**Files Modified:**
- `frontend/src/components/tools/ToolCard.vue` (add compare interaction)

---

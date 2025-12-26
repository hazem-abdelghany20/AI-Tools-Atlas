### Story 2.7: Navigate to Tool Profile

**As a** user,
**I want** to click on a tool card and view its detailed profile,
**So that** I can learn more about a specific tool.

**Acceptance Criteria:**

**Given** I see a tool card in search results or category view
**When** I click on the card (anywhere except bookmark/compare buttons)
**Then** I navigate to `/tools/:slug` and the ToolProfileView loads

**And:**
- Clicking the tool name also navigates to the profile
- Clicking the logo navigates to the profile
- Bookmark and compare buttons do NOT trigger navigation (event.stopPropagation)

**Router Configuration:**
- Route defined: `{ path: '/tools/:slug', component: ToolProfileView }`
- Route receives slug as param: `route.params.slug`

**ToolProfileView (Basic Shell):**
- For this story, create a minimal profile view showing:
  - Tool name as heading
  - "Loading full profile..." message or skeleton
  - This will be fully implemented in Epic 3
- Verifies navigation works end-to-end

**Technical Implementation:**

- ToolCard component:
  ```vue
  <div class="tool-card" @click="navigateToProfile">
    <!-- card content -->
    <button @click.stop="toggleBookmark">Bookmark</button>
    <button @click.stop="addToCompare">Compare</button>
  </div>
  ```
- navigateToProfile method: `router.push(\`/tools/${props.tool.slug}\`)`
- ToolProfileView.vue stub: `<div>Tool Profile: {{ route.params.slug }}</div>`

**Prerequisites:** Story 2.4 (ToolCard component)

**Files Created:**
- `frontend/src/views/ToolProfileView.vue` (basic shell for now, will be expanded in Epic 3)

---

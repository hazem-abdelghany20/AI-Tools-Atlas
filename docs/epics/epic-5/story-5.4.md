### Story 5.4: Comparison View

**As a** user,
**I want** to see tools side-by-side in a comparison table,
**So that** I can evaluate options and make a decision.

**Acceptance Criteria:**

**CompareView (`/compare?tools=slug1,slug2,slug3`):**
- Page title: "المقارنة" (Comparison)
- Side-by-side table with tools as columns (2-4 tools)
- Comparison rows:
  - **Overview:** Name, logo, "best for", tagline
  - **Category:** Primary category
  - **Pricing:** Free tier, starting price, billing model
  - **Rating:** Overall rating, review count
  - **Features:** Key features (if available)
  - **Social Proof:** Bookmark count, badges
  - **Actions:** "Visit Tool" button, bookmark button for each tool
- "Remove from comparison" button on each column
- "Share comparison" button: copies URL to clipboard
- Empty state: Prompt to add tools from search or bookmarks

**Responsive:**
- Desktop: full table with all columns visible
- Mobile: horizontal scroll, sticky first column (row labels)

**Technical Implementation:**
- Parse tool slugs from URL: `route.query.tools.split(',')`
- Fetch all tools: `Promise.all(slugs.map(slug => apiClient.get(\`/tools/${slug}\`)))`
- CompareTable.vue component builds table from tool data
- Remove tool: update URL, refetch remaining tools

**Prerequisites:** Story 5.3 (Comparison selection), Epic 3 (Tool data)

**Files Created:**
- `frontend/src/views/CompareView.vue`
- `frontend/src/components/compare/CompareTable.vue`

---

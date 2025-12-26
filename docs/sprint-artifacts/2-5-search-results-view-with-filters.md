### Story 2.5: Search Results View with Filters

**As a** user,
**I want** to see search results with filtering and sorting options,
**So that** I can narrow down to the most relevant tools for my needs.

**Acceptance Criteria:**

**Given** I navigate to `/search?q=my+query` or `/categories/:slug`
**When** the SearchResultsView loads
**Then** I see:

**Page Layout:**
- **Filter panel (right side in RTL, drawer on mobile):**
  - Categories filter: dropdown or list of checkboxes
  - Price filter: radio buttons (الكل، مجاني، Freemium، مدفوع)
  - Rating filter: slider or radio (All, 4+, 4.5+)
  - Platform filter: checkboxes (Web, Mobile, API, Desktop)
  - "مسح الفلاتر" (Clear filters) button
- **Results area (main content):**
  - Header showing: query/category name, result count, sort dropdown
  - Sort dropdown: "الأعلى تقييماً" (Top rated), "الأكثر حفظاً" (Most bookmarked), "الأحدث" (Newest), "الرائج" (Trending)
  - Grid of ToolCard components (3-4 columns desktop, 1 column mobile)
  - Pagination controls at bottom (Previous/Next, page numbers)
  - Empty state if no results: "لم تجد ما تريد؟ استعرض حسب الفئة" with links to categories

**Filter Behavior:**
- Changing any filter updates URL query params: `?category=...&price=...&min_rating=...&sort=...`
- URL params drive the API call and filter UI state
- Page resets to 1 when filters change
- Filters are visible/sticky on desktop, drawer on mobile with "تصفية" (Filter) button

**Loading States:**
- Skeleton loaders for tool cards while loading (show 8-12 skeletons)
- No blocking spinner, results area shows skeletons in grid layout
- Preserve previous results briefly while new results load (optional optimization)

**Search Query Handling:**
- Free-text query from hero search calls `GET /api/v1/search/tools?q=...`
- Category browsing calls `GET /api/v1/categories/:slug/tools`
- All other filters apply to both search and category views

**And:**
- Uses `useFiltersStore` to sync filter state with URL params
- Filter changes trigger new API call: `apiClient.get('/tools', { params: filters })`
- Pagination: `&page=2` in query params

**Technical Implementation:**

- **SearchResultsView.vue** (`frontend/src/views/SearchResultsView.vue`)
- Watch route query params, fetch tools when params change:
  ```typescript
  watch(() => route.query, () => {
    fetchTools();
  }, { immediate: true });
  ```
- Filter panel: separate component `FiltersPanel.vue` in `components/search/`
- Responsive filter drawer using Headless UI or custom modal on mobile
- URL sync: use Vue Router to update query params when filters change
- Pagination: `<Pagination :current-page="page" :total-pages="totalPages" @page-change="goToPage" />`

**Prerequisites:** Story 2.2 (Tools API), Story 2.4 (ToolCard component)

**Files Created:**
- `frontend/src/views/SearchResultsView.vue`
- `frontend/src/components/search/FiltersPanel.vue`
- `frontend/src/components/common/Pagination.vue`

---

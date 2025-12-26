### Story 3.3: Tool Profile View - Hero & Overview

**As a** user,
**I want** to see a rich tool profile with hero section and overview,
**So that** I can quickly understand what the tool is and if it fits my needs.

**Acceptance Criteria:**

**Given** I navigate to `/tools/:slug`
**When** the ToolProfileView loads
**Then** I see the following sections:

**Hero Section:**
- Tool logo (large, right side in RTL) with fallback if missing
- Tool name as main heading (H1)
- Tagline as subheading
- Overall rating (stars + numeric, e.g., "4.5") with review count "(42 مراجعة)"
- Pricing badge: "مجاني" / "Freemium" / starting price
- Primary CTA button: "زيارة الأداة" (Visit Tool) - opens official_url in new tab
- Bookmark button (larger than card version): "حفظ" / "محفوظ"
- Category pill/link: navigates to category page
- Badges displayed as pills: "الأفضل في الفئة", "نجم صاعد", etc.

**Overview Section:**
- **"الأفضل لـ" (Best for)** - displayed prominently, 2-3 line summary
- **Description** - full tool description, multiple paragraphs
- **Primary Use Cases** - bullet list or pills showing main use cases
- **Target Roles** - who this tool is for (e.g., "المحامون، الكتّاب، المسوقون")
- **Platforms** - icons or text showing Web, Mobile, API, Desktop availability

**Loading State:**
- Skeleton loaders for hero and overview while API call is in progress
- No blocking spinner, layout structure visible with skeleton content

**Error Handling:**
- 404 page if tool not found: "الأداة غير موجودة" with link back to home
- Network error: retry button

**UX Requirements:**
- Dark background, neon blue for CTA and active states
- RTL-aware layout with proper Arabic typography (Cairo font)
- Responsive: single column on mobile, optimized spacing on desktop
- "Visit Tool" button is visually prominent (neon blue, larger size)

**Technical Implementation:**

- Fetch tool data on component mount: `apiClient.get(\`/tools/${route.params.slug}\`)`
- Use Vue reactive refs for tool data, loading, error states
- Bookmark interaction: `useBookmarksStore` to add/remove bookmark
- External link opens in new tab: `<a :href="tool.official_url" target="_blank" rel="noopener">`

**Prerequisites:** Story 3.1 (Tool profile API), Epic 1 (Frontend components)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (expand from Story 2.7 shell)

**New Files:**
- `frontend/src/components/tools/ToolHero.vue` (optional sub-component)
- `frontend/src/components/tools/ToolOverview.vue` (optional sub-component)

---

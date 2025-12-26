### Story 2.4: Tool Card Component

**As a** developer,
**I want** a reusable ToolCard component,
**So that** tools can be displayed consistently across search results, categories, and bookmarks.

**Acceptance Criteria:**

**Given** I need to display tool information in card format
**When** I create the ToolCard component
**Then** it displays the following:

**Card Layout (matches G2-style from UX Design):**
- Tool logo (left in RTL) with fallback placeholder if missing
- Tool name as heading
- "Best for" line in lighter text (one-line truncate)
- Rating display: stars (or numeric) + review count in parentheses "(23 مراجعة)"
- Pricing badge: "مجاني" (Free), "Freemium", or starting price
- Tags: 2-3 primary tags as small pills (overflow hidden)
- Bookmark button (hollow heart/bookmark icon, filled when bookmarked)
- "مقارنة" (Compare) button/checkbox to add to comparison
- Clicking card (except buttons) navigates to `/tools/:slug`

**Visual Design (from UX Design spec):**
- Dark card background (#0A0B10) with subtle border (#1F2933)
- Hover state: slight border glow (neon blue)
- RTL-aware: logo on right, actions on left
- Responsive: full-width on mobile, fixed width/flex-basis on desktop grid
- Cairo font for all text
- Proper spacing (8px base unit from Architecture)

**Component Props:**
```typescript
interface ToolCardProps {
  tool: {
    id: number;
    slug: string;
    name: string;
    logo_url?: string;
    best_for: string;
    avg_rating_overall: number;
    review_count: number;
    pricing_summary: string;
    has_free_tier: boolean;
    tags: Array<{ name: string }>;
  };
  showCompare?: boolean; // default true
}
```

**Bookmark Interaction:**
- Uses `useBookmarksStore` to check `isBookmarked(tool.id)`
- Click bookmark button: calls `addBookmark(tool.id)` or `removeBookmark(tool.id)`
- Optimistic UI: immediately updates icon state, rollback on error
- Shows toast on error: "فشل في حفظ الأداة" (Failed to save tool)

**Compare Interaction:**
- Click compare button: emits `@add-to-compare` event with tool
- Parent component handles comparison state (will be in ComparisonStore later)

**Technical Implementation:**

- Vue 3 SFC in `frontend/src/components/tools/ToolCard.vue`
- Use Tailwind for styling, no scoped CSS unless necessary
- Rating display: create small `RatingStars.vue` component or use numeric display
- Icons: use a lightweight icon library or inline SVGs for bookmark/compare
- Truncate "best for" text: `class="truncate"` or `-webkit-line-clamp: 1`

**Prerequisites:** Epic 1 (Frontend setup), Story 2.2 (Tools API for data structure)

**Files Created:**
- `frontend/src/components/tools/ToolCard.vue`
- `frontend/src/components/common/RatingStars.vue` (optional)

---

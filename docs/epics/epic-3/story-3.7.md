### Story 3.7: Social Proof & Engagement Indicators

**As a** user,
**I want** to see social proof indicators like badges and bookmark counts,
**So that** I can gauge the tool's popularity and credibility.

**Acceptance Criteria:**

**Given** I am viewing a tool profile
**When** the page loads
**Then** I see social proof indicators:

**Badges (already in Story 3.3 hero):**
- Displayed prominently in hero section as pills/chips
- Examples: "الأفضل في الفئة" (Top in Category), "نجم صاعد" (Rising Star), "اختيار المحرر" (Editor's Pick)
- Each badge has icon (if available) + name
- Color-coded or styled distinctively from regular tags

**Bookmark Count:**
- Displayed near bookmark button or in stats area
- Format: "حفظه 234 مستخدماً" (Bookmarked by 234 users)
- Updates when user bookmarks/unbookmarks (optimistic UI)

**Review Count & Rating:**
- Already covered in hero section (Story 3.3)
- Overall rating with stars + numeric value
- Review count as clickable link scrolling to reviews section (will be in Epic 4)

**Engagement Stats (Optional):**
- Small stats panel showing:
  - 🔖 Bookmark count
  - ⭐ Rating (avg_rating_overall)
  - 💬 Review count
- Displayed below hero or in sidebar (desktop)

**Visual Design:**
- Badges use neon blue or distinct colors, slightly elevated (subtle shadow/glow)
- Stats displayed with icons and clear Arabic labels
- RTL-aware layout for badge row
- Subtle animations on hover for badges

**Technical Implementation:**

- Badges already loaded in tool profile API (Story 3.1)
- Display: iterate over `tool.badges` array
- Bookmark count from `tool.bookmark_count`
- Update bookmark count optimistically when user bookmarks:
  ```typescript
  const displayBookmarkCount = computed(() => {
    return tool.value.bookmark_count + (isBookmarked ? 1 : 0);
  });
  ```

**Prerequisites:** Story 3.3 (Hero section), Story 3.1 (Tool data with badges)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (ensure badges and stats are displayed)
- `frontend/src/components/tools/ToolHero.vue` (if using sub-component)

---

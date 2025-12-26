### Story 8.4: Moderation Queue UI ✅ DONE

**As a** moderator,
**I want** a moderation queue UI,
**So that** I can review and act on reported content.

**Acceptance Criteria:**

**Moderation Queue Page (`/moderation/queue`):**
- Table of reports with columns: Type (Tool/Review), Content Preview, Reason, Reporter, Date, Status, Actions
- Filters: Type (All/Tools/Reviews), Status (Pending/Reviewed/Dismissed)
- Pagination

**Report Details:**
- Click report row expands to show:
  - Full content (tool details or review text)
  - Reporter info (if available)
  - Report reason and comment
  - "View in context" link (opens tool profile or review in new tab)

**Actions:**
- For reviews: Approve, Hide, Remove buttons
- For tools: View on site, Flag for admin review (tools handled by admins in Story 7.2)
- Dismiss report button (marks report as reviewed without action)
- Confirmation modals for destructive actions

**Moderation History:**
- Link to view history for a review
- Shows timeline of actions: who approved/hid/removed and when

**Visual Design:**
- Clear status indicators: Pending (yellow), Reviewed (green), Dismissed (gray)
- Content preview truncated with "Show more"
- Actions clearly labeled

**Technical Implementation:**
- ModerationQueueView.vue fetches reports
- Action buttons call moderation endpoints
- Optimistic UI: update status immediately, rollback on error
- Role check: only moderators/admins can access

**Prerequisites:** Story 8.3

**New Files:**
- `frontend/src/views/moderation/ModerationQueueView.vue`
- `frontend/src/router/moderation.ts` (moderation routes)

---

## Change Log:
| Date | Change |
|------|--------|
| 2025-12-26 | Code Review: MEDIUM - Fixed formatDate validation, added Escape key handlers for modals |

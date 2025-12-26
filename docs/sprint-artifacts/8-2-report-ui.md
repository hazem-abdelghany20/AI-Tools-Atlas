### Story 8.2: Report UI

**As a** user,
**I want** to report tools or reviews,
**So that** I can flag inappropriate content.

**Acceptance Criteria:**

**Report Button:**
- On tool profiles: "إبلاغ" (Report) link near tool name
- On review cards: "إبلاغ" (Report) link in review card footer

**Report Modal:**
- Dropdown: select reason (Spam, Abuse, Misinformation, Other)
- Textarea: optional comment (max 500 chars)
- Submit button: "إرسال البلاغ"
- Cancel button

**Submission:**
- On submit: POST to report endpoint
- Success: close modal, show toast "شكراً للإبلاغ. سنراجعه قريباً"
- Error: show error message in modal

**Technical Implementation:**
- ReportModal.vue component
- Accept props: reportableType, reportableId
- API call to appropriate report endpoint

**Prerequisites:** Story 8.1

**New Files:**
- `frontend/src/components/moderation/ReportModal.vue`

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (add report button)
- `frontend/src/components/reviews/ReviewCard.vue` (add report button)

---

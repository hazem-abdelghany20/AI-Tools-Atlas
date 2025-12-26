### Story 4.4: Aggregated Ratings Display

**As a** user,
**I want** to see aggregated ratings at the top of the reviews section,
**So that** I can quickly understand overall sentiment.

**Acceptance Criteria:**

**Ratings Summary Component:**
- Overall rating (large): 4.5 / 5 stars with total review count
- Rating distribution: bar chart showing count per star (5★: 42, 4★: 18, etc.)
- Dimension ratings: horizontal bars for Ease of use, Value, Accuracy, Speed, Support
- Displayed above the reviews list
- Updates when new reviews are submitted

**Visual Design:**
- Clean, scannable layout
- Neon blue for filled stars and bars
- Rating distribution bars clearly show proportion
- Arabic labels for all dimensions

**Technical Implementation:**
- Data from tool object: avg_rating_overall, review_count
- Dimension ratings: calculate averages from reviews or store on tool table
- Rating distribution: fetch from backend or calculate client-side

**Prerequisites:** Story 4.1, Story 4.2

**New Files:**
- `frontend/src/components/reviews/RatingsAggregation.vue`

---

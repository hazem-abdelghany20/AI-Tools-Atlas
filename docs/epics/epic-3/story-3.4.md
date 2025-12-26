### Story 3.4: Tool Profile View - Features & Pricing

**As a** user,
**I want** to see detailed features and pricing information,
**So that** I can understand what the tool offers and how much it costs.

**Acceptance Criteria:**

**Given** I am viewing a tool profile
**When** I scroll below the overview
**Then** I see:

**Features Section:**
- **Heading:** "المميزات" (Features)
- List or grid of key features
- Each feature can have:
  - Icon (optional)
  - Feature name/title
  - Short description (optional)
- Displayed as cards or bulleted list depending on data structure
- If no features data available, section is hidden or shows placeholder

**Pricing Section:**
- **Heading:** "الأسعار" (Pricing)
- Free tier information (if `has_free_tier` is true):
  - "✓ يتوفر خطة مجانية" (Free tier available)
- Pricing summary displayed clearly (from `pricing_summary` field):
  - E.g., "يبدأ من $10/شهرياً" (Starts at $10/month)
  - Or structured tiers if data supports it (Basic, Pro, Enterprise)
- Link to official pricing page: "عرض تفاصيل الأسعار" (View pricing details) linking to official_url
- If no pricing data, show: "الرجاء زيارة الموقع الرسمي للأسعار" (Please visit official site for pricing)

**Visual Design:**
- Sections use consistent heading style (H2, Cairo font, white text)
- Features displayed in 2-3 column grid on desktop, stacked on mobile
- Pricing displayed as structured cards or simple text block
- Proper spacing between sections (24-32px)

**Technical Implementation:**

- Parse `tool.pricing_summary` and display as formatted text
- Features: if stored as JSON array, iterate and display; if text, display as-is
- Conditional rendering: `v-if="tool.features && tool.features.length > 0"`
- Pricing tiers: could be hardcoded structure or parsed from pricing_summary

**Prerequisites:** Story 3.3 (Tool profile view structure)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (add sections)

**New Files (optional):**
- `frontend/src/components/tools/ToolFeatures.vue`
- `frontend/src/components/tools/ToolPricing.vue`

---

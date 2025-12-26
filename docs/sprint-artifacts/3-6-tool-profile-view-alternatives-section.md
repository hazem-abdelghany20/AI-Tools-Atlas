### Story 3.6: Tool Profile View - Alternatives Section

**As a** user,
**I want** to see similar and alternative tools on the profile,
**So that** I can explore other options before making a decision.

**Acceptance Criteria:**

**Given** I am viewing a tool profile
**When** the alternatives section loads
**Then** I see:

**Similar Tools Subsection:**
- **Heading:** "أدوات مشابهة" (Similar Tools)
- Horizontal scrolling row of tool cards (using ToolCard component)
- Shows 4-6 similar tools
- Cards are clickable, navigate to respective tool profiles
- If no similar tools: section is hidden

**Alternatives Subsection:**
- **Heading:** "بدائل لهذه الأداة" (Alternatives to this tool)
- Horizontal scrolling row of tool cards
- Shows 4-6 alternative tools
- Cards are clickable, navigate to respective tool profiles
- If no alternatives: section is hidden

**Visual Design:**
- Each subsection clearly separated with heading
- Horizontal scroll with subtle scroll indicators (gradient fade on edges)
- Tool cards use same ToolCard component as search results (consistent design)
- RTL-aware scrolling: scroll starts from right side
- Mobile: scroll snap for smooth scrolling experience

**Technical Implementation:**

- Fetch alternatives on component mount: `apiClient.get(\`/tools/${route.params.slug}/alternatives\`)`
- Store in reactive refs: `similar`, `alternatives`
- Conditional rendering: only show sections if arrays have items
- Horizontal scroll container:
  ```vue
  <div class="flex overflow-x-auto gap-4 pb-4">
    <ToolCard v-for="tool in similar" :key="tool.id" :tool="tool" />
  </div>
  ```
- RTL scrolling: use `dir="rtl"` on scroll container or CSS `direction: rtl`

**Prerequisites:** Story 3.2 (Alternatives API), Story 2.4 (ToolCard component)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (add alternatives section)

---

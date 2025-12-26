### Story 3.5: Tool Profile View - Media Gallery

**As a** user,
**I want** to see screenshots and videos of the tool,
**So that** I can visually understand how it works.

**Acceptance Criteria:**

**Given** I am viewing a tool profile
**When** the media section loads
**Then** I see:

**Media Section:**
- **Heading:** "صور وفيديوهات" (Screenshots and Videos)
- Media items displayed in order (from `display_order` field)
- **Screenshots:**
  - Displayed as thumbnail grid (2-3 columns on desktop, 1-2 on mobile)
  - Click thumbnail opens lightbox/modal with full-size image
  - Lightbox has prev/next navigation for multiple images
  - Image alt text for accessibility
- **Videos (YouTube embeds):**
  - Displayed as embedded YouTube player or thumbnail
  - Lazy-loaded: only load iframe when user scrolls to media section or clicks play
  - Responsive embed: maintains 16:9 aspect ratio
  - If video URL is YouTube, extract video ID and embed properly
- **Mixed media:**
  - Screenshots displayed first, then videos
  - Or single gallery with media type indicators

**Empty State:**
- If no media: section is hidden or shows "لا توجد صور أو فيديوهات متاحة" (No media available)

**Performance:**
- Lazy load YouTube iframes to avoid blocking page render (critical from UX Design spec)
- Use `loading="lazy"` for screenshot images
- Skeleton placeholders while media loads

**Technical Implementation:**

- **ToolMediaGallery.vue** component
- Iterate over `tool.media` array:
  ```vue
  <div v-for="media in tool.media" :key="media.id">
    <img v-if="media.type === 'screenshot'"
         :src="media.thumbnail_url || media.url"
         @click="openLightbox(media)" />
    <YouTubeEmbed v-else-if="media.type === 'video'"
                  :url="media.url"
                  lazy />
  </div>
  ```
- Lightbox modal: use a lightweight library or custom Vue component
- YouTube lazy loading: render thumbnail with play button, only load iframe on click

**Prerequisites:** Story 3.3 (Tool profile structure)

**Files Modified:**
- `frontend/src/views/ToolProfileView.vue` (add media section)

**New Files:**
- `frontend/src/components/tools/ToolMediaGallery.vue`
- `frontend/src/components/common/Lightbox.vue` (for image viewing)
- `frontend/src/components/common/YouTubeEmbed.vue` (lazy-loaded YouTube player)

---

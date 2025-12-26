<template>
  <div class="tool-stats flex flex-wrap gap-4 md:gap-6">
    <!-- Rating Stat -->
    <div class="stat-item flex items-center gap-2">
      <div class="stat-icon w-10 h-10 rounded-lg bg-yellow-500/20 flex items-center justify-center">
        <svg class="w-5 h-5 text-yellow-400" fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
      </div>
      <div>
        <div class="text-white font-bold">{{ rating }}</div>
        <div class="text-gray-400 text-xs">التقييم</div>
      </div>
    </div>

    <!-- Review Count Stat -->
    <div class="stat-item flex items-center gap-2">
      <div class="stat-icon w-10 h-10 rounded-lg bg-blue-500/20 flex items-center justify-center">
        <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      </div>
      <div>
        <div class="text-white font-bold">{{ reviewCount }}</div>
        <div class="text-gray-400 text-xs">مراجعة</div>
      </div>
    </div>

    <!-- Bookmark Count Stat -->
    <div class="stat-item flex items-center gap-2">
      <div class="stat-icon w-10 h-10 rounded-lg bg-red-500/20 flex items-center justify-center">
        <svg class="w-5 h-5 text-red-400" fill="currentColor" viewBox="0 0 24 24">
          <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
        </svg>
      </div>
      <div>
        <div class="text-white font-bold">{{ displayBookmarkCount }}</div>
        <div class="text-gray-400 text-xs">حفظ</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Tool } from '../../lib/types'
import { useBookmarksStore } from '../../stores/bookmarks'

const props = defineProps<{
  tool: Tool
}>()

const bookmarksStore = useBookmarksStore()

const rating = computed(() => {
  return props.tool.avg_rating_overall?.toFixed(1) || '0.0'
})

const reviewCount = computed(() => {
  return props.tool.review_count || 0
})

// Optimistic bookmark count: adjust based on user's action
const displayBookmarkCount = computed(() => {
  const baseCount = props.tool.bookmark_count || 0
  const isBookmarked = bookmarksStore.isBookmarked(props.tool.id)
  // If user bookmarked and it wasn't already counted, add 1
  // This is optimistic UI - the real count comes from API
  return baseCount + (isBookmarked ? 1 : 0)
})
</script>

<style scoped>
.stat-item {
  transition: transform 0.2s ease;
}

.stat-item:hover {
  transform: translateY(-2px);
}
</style>

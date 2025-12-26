<template>
  <article
    class="tool-card group relative bg-dark-surface border border-gray-800 rounded-xl p-4
           hover:border-neon-blue/50 transition-all duration-200 cursor-pointer"
    dir="rtl"
    @click="handleCardClick"
  >
    <!-- Main Content -->
    <div class="flex gap-4">
      <!-- Tool Logo -->
      <div class="flex-shrink-0">
        <img
          v-if="tool.logo_url"
          :src="tool.logo_url"
          :alt="tool.name"
          class="w-16 h-16 rounded-xl object-cover bg-gray-800"
        />
        <div
          v-else
          class="w-16 h-16 rounded-xl bg-neon-blue/20 flex items-center justify-center"
        >
          <span class="text-neon-blue font-bold text-2xl">
            {{ tool.name.charAt(0) }}
          </span>
        </div>
      </div>

      <!-- Tool Info -->
      <div class="flex-1 min-w-0">
        <!-- Name & Pricing Badge -->
        <div class="flex items-start justify-between gap-2">
          <h3 class="font-semibold text-white text-lg truncate">
            {{ tool.name }}
          </h3>
          <span
            class="flex-shrink-0 px-2 py-0.5 text-xs font-medium rounded-full"
            :class="pricingBadgeClass"
          >
            {{ pricingLabel }}
          </span>
        </div>

        <!-- Best For -->
        <p class="text-gray-400 text-sm mt-1 truncate">
          {{ tool.best_for || tool.tagline || 'أداة ذكاء اصطناعي' }}
        </p>

        <!-- Rating -->
        <div class="mt-2">
          <RatingStars
            :rating="tool.avg_rating_overall"
            :review-count="tool.review_count"
          />
        </div>

        <!-- Tags (hidden in compact mode) -->
        <div v-if="!compact && tool.tags && tool.tags.length > 0" class="flex flex-wrap gap-1 mt-3">
          <span
            v-for="(tag, index) in displayedTags"
            :key="tag.id || index"
            class="px-2 py-0.5 text-xs bg-gray-800 text-gray-400 rounded-full"
          >
            {{ tag.name }}
          </span>
          <span
            v-if="tool.tags.length > maxTags"
            class="px-2 py-0.5 text-xs bg-gray-800 text-gray-500 rounded-full"
          >
            +{{ tool.tags.length - maxTags }}
          </span>
        </div>
      </div>
    </div>

    <!-- Actions (simplified in compact mode) -->
    <div class="flex items-center justify-between mt-4 pt-3 border-t border-gray-800" :class="{ 'mt-2 pt-2': compact }">
      <div class="flex items-center gap-2">
        <!-- Bookmark Button -->
        <button
          @click.stop="toggleBookmark"
          class="p-2 rounded-lg hover:bg-gray-800 transition-colors"
          :class="{ 'text-red-500': isBookmarkedLocal, 'text-gray-400': !isBookmarkedLocal }"
          :title="isBookmarkedLocal ? 'إزالة من المحفوظات' : 'حفظ'"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="w-5 h-5"
            :fill="isBookmarkedLocal ? 'currentColor' : 'none'"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
            />
          </svg>
        </button>

        <!-- Compare Button -->
        <button
          v-if="showCompare"
          @click.stop="handleCompare"
          class="p-2 rounded-lg hover:bg-gray-800 transition-colors"
          :class="{ 'text-neon-blue': isInCompare, 'text-gray-400': !isInCompare }"
          title="مقارنة"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="w-5 h-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
            />
          </svg>
        </button>
      </div>

      <!-- View Details Link -->
      <span class="text-neon-blue text-sm group-hover:underline">
        عرض التفاصيل
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="inline w-4 h-4 mr-1"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M15 19l-7-7 7-7"
          />
        </svg>
      </span>
    </div>
  </article>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useBookmarksStore } from '@/stores/bookmarks'
import RatingStars from '@/components/common/RatingStars.vue'
import type { Tool } from '@/lib/types'

const props = withDefaults(defineProps<{
  tool: Tool
  showCompare?: boolean
  isInCompare?: boolean
  compact?: boolean
}>(), {
  showCompare: true,
  isInCompare: false,
  compact: false
})

const emit = defineEmits<{
  (e: 'add-to-compare', tool: Tool): void
  (e: 'remove-from-compare', tool: Tool): void
}>()

const router = useRouter()
const bookmarksStore = useBookmarksStore()

const maxTags = 3
const isBookmarkedLocal = ref(false)

// Initialize local bookmark state from store
watch(
  () => bookmarksStore.isBookmarked(props.tool.id),
  (newVal) => {
    isBookmarkedLocal.value = newVal
  },
  { immediate: true }
)

const displayedTags = computed(() => {
  if (!props.tool.tags) return []
  return props.tool.tags.slice(0, maxTags)
})

const pricingLabel = computed(() => {
  if (props.tool.has_free_tier) {
    if (props.tool.pricing_summary?.toLowerCase().includes('free')) {
      return 'مجاني'
    }
    return 'Freemium'
  }
  return props.tool.pricing_summary || 'مدفوع'
})

const pricingBadgeClass = computed(() => {
  if (props.tool.has_free_tier) {
    return 'bg-green-500/20 text-green-400'
  }
  return 'bg-gray-700 text-gray-300'
})

function handleCardClick() {
  router.push(`/tools/${props.tool.slug}`)
}

async function toggleBookmark() {
  const wasBookmarked = isBookmarkedLocal.value
  isBookmarkedLocal.value = !wasBookmarked // Optimistic update

  try {
    if (wasBookmarked) {
      await bookmarksStore.removeBookmark(props.tool.id)
    } else {
      await bookmarksStore.addBookmark(props.tool.id)
    }
  } catch (error) {
    // Rollback on error
    isBookmarkedLocal.value = wasBookmarked
    console.error('Failed to toggle bookmark:', error)
  }
}

function handleCompare() {
  if (props.isInCompare) {
    emit('remove-from-compare', props.tool)
  } else {
    emit('add-to-compare', props.tool)
  }
}
</script>

<style scoped>
.tool-card {
  min-height: 180px;
}
</style>

<template>
  <section class="tool-hero bg-dark-surface rounded-lg p-6 md:p-8">
    <div class="flex flex-col md:flex-row-reverse gap-6 items-start">
      <!-- Logo -->
      <div class="flex-shrink-0">
        <div
          v-if="tool.logo_url"
          class="w-24 h-24 md:w-32 md:h-32 rounded-xl overflow-hidden bg-dark-border"
        >
          <img
            :src="tool.logo_url"
            :alt="tool.name"
            class="w-full h-full object-cover"
          />
        </div>
        <div
          v-else
          class="w-24 h-24 md:w-32 md:h-32 rounded-xl bg-dark-border flex items-center justify-center"
        >
          <span class="text-4xl md:text-5xl text-gray-400">{{ tool.name?.charAt(0) || '?' }}</span>
        </div>
      </div>

      <!-- Main Content -->
      <div class="flex-1">
        <!-- Category & Badges -->
        <div class="flex flex-wrap items-center gap-2 mb-3">
          <router-link
            v-if="tool.primary_category"
            :to="`/categories/${tool.primary_category.slug}`"
            class="px-3 py-1 bg-dark-border text-sm rounded-full hover:bg-dark-border/80 transition-colors"
          >
            {{ tool.primary_category.name }}
          </router-link>
          <span
            v-for="badge in tool.badges"
            :key="badge.id"
            class="px-3 py-1 bg-neon-blue/20 text-neon-blue text-sm rounded-full"
          >
            {{ badge.name }}
          </span>
        </div>

        <!-- Name & Tagline -->
        <h1 class="text-2xl md:text-4xl font-bold mb-2">{{ tool.name }}</h1>
        <p v-if="tool.tagline" class="text-gray-400 text-lg mb-4">{{ tool.tagline }}</p>

        <!-- Rating & Pricing -->
        <div class="flex flex-wrap items-center gap-4 mb-6">
          <div class="flex items-center gap-2">
            <RatingStars :rating="tool.avg_rating_overall" />
            <span class="text-white font-semibold">{{ tool.avg_rating_overall?.toFixed(1) || '0.0' }}</span>
            <span class="text-gray-400">({{ tool.review_count || 0 }} مراجعة)</span>
          </div>
          <span class="pricing-badge px-3 py-1 rounded-full text-sm font-medium" :class="pricingClass">
            {{ pricingLabel }}
          </span>
        </div>

        <!-- CTA Buttons -->
        <div class="flex flex-wrap items-center gap-3">
          <a
            v-if="tool.official_url"
            :href="tool.official_url"
            target="_blank"
            rel="noopener noreferrer"
            class="cta-button px-6 py-3 bg-neon-blue text-dark-bg font-bold rounded-lg hover:bg-neon-blue-hover transition-colors inline-flex items-center gap-2"
          >
            <span>زيارة الأداة</span>
            <svg class="w-5 h-5 rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
          </a>
          <button
            @click="toggleBookmark"
            class="bookmark-button px-6 py-3 rounded-lg font-medium transition-colors inline-flex items-center gap-2"
            :class="bookmarked ? 'bg-neon-blue text-dark-bg' : 'bg-dark-border text-white hover:bg-dark-border/80'"
          >
            <svg class="w-5 h-5" :fill="bookmarked ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
            </svg>
            <span>{{ bookmarked ? 'محفوظ' : 'حفظ' }}</span>
          </button>
          <button
            @click="$emit('report')"
            class="px-4 py-3 rounded-lg text-gray-400 hover:text-red-400 hover:bg-dark-border transition-colors inline-flex items-center gap-2"
            title="الإبلاغ عن هذه الأداة"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <span class="text-sm">إبلاغ</span>
          </button>
        </div>

        <!-- Bookmark Count -->
        <div v-if="displayBookmarkCount > 0" class="mt-4 text-gray-400 text-sm">
          <svg class="w-4 h-4 inline-block me-1" fill="currentColor" viewBox="0 0 24 24">
            <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
          </svg>
          حفظه {{ displayBookmarkCount }} مستخدماً
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Tool } from '../../lib/types'
import RatingStars from '../common/RatingStars.vue'
import { useBookmarksStore } from '../../stores/bookmarks'

const props = defineProps<{
  tool: Tool
}>()

const emit = defineEmits<{
  (e: 'bookmark-toggle', toolId: number): void
  (e: 'report'): void
}>()

const bookmarksStore = useBookmarksStore()

const bookmarked = computed(() => bookmarksStore.isBookmarked(props.tool.id))

const pricingLabel = computed(() => {
  if (props.tool.has_free_tier && !props.tool.pricing_summary) {
    return 'مجاني'
  }
  if (props.tool.has_free_tier) {
    return 'Freemium'
  }
  return props.tool.pricing_summary || 'مدفوع'
})

const pricingClass = computed(() => {
  if (props.tool.has_free_tier && !props.tool.pricing_summary) {
    return 'bg-green-500/20 text-green-400'
  }
  if (props.tool.has_free_tier) {
    return 'bg-yellow-500/20 text-yellow-400'
  }
  return 'bg-blue-500/20 text-blue-400'
})

// Optimistic bookmark count display
const displayBookmarkCount = computed(() => {
  const baseCount = props.tool.bookmark_count || 0
  // Adjust count based on user's bookmark action
  return baseCount + (bookmarked.value ? 1 : 0)
})

const toggleBookmark = () => {
  emit('bookmark-toggle', props.tool.id)
}
</script>

<style scoped>
.cta-button {
  box-shadow: 0 0 20px rgba(0, 212, 255, 0.3);
}

.cta-button:hover {
  box-shadow: 0 0 30px rgba(0, 212, 255, 0.5);
}
</style>

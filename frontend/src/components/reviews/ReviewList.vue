<template>
  <section class="reviews-section bg-dark-surface rounded-xl p-6" dir="rtl">
    <!-- Ratings Aggregation (shown when we have reviews) -->
    <RatingsAggregation
      v-if="reviews.length > 0 && !loading"
      :avg-rating="avgRating"
      :review-count="totalReviews"
      :reviews="reviews"
    />

    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <h2 class="text-xl font-bold text-white">المراجعات</h2>
        <span v-if="totalReviews > 0" class="text-gray-400 text-sm">
          ({{ totalReviews }})
        </span>
      </div>

      <!-- Sort Dropdown -->
      <div v-if="reviews.length > 0" class="relative">
        <select
          v-model="currentSort"
          class="bg-dark-bg border border-gray-700 rounded-lg px-4 py-2 text-sm text-gray-300
                 focus:border-neon-blue focus:outline-none appearance-none pr-8 cursor-pointer"
          @change="handleSortChange"
        >
          <option value="newest">الأحدث</option>
          <option value="most_helpful">الأكثر فائدة</option>
          <option value="highest">الأعلى تقييماً</option>
          <option value="lowest">الأدنى تقييماً</option>
        </select>
        <svg
          class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="space-y-4">
      <div v-for="n in 3" :key="n" class="animate-pulse">
        <div class="bg-dark-bg rounded-xl p-5 space-y-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-gray-700"></div>
            <div class="space-y-2">
              <div class="h-4 bg-gray-700 rounded w-24"></div>
              <div class="h-3 bg-gray-700 rounded w-32"></div>
            </div>
          </div>
          <div class="h-16 bg-gray-700 rounded"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="reviews.length === 0" class="text-center py-12">
      <div class="text-4xl mb-4">
        <svg class="w-16 h-16 mx-auto text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
      </div>
      <p class="text-gray-400 mb-4">لا توجد مراجعات بعد. كن أول من يراجع!</p>
      <button
        v-if="showAddButton"
        class="px-6 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg
               hover:bg-neon-blue-hover transition-colors"
        @click="$emit('add-review')"
      >
        إضافة مراجعة
      </button>
    </div>

    <!-- Reviews List -->
    <div v-else class="space-y-4">
      <ReviewCard
        v-for="review in reviews"
        :key="review.id"
        :review="review"
        @helpful="handleHelpful"
        @report="handleReport"
      />

      <!-- Load More Button -->
      <div v-if="hasMore" class="text-center pt-4">
        <button
          class="px-6 py-2 border border-gray-700 text-gray-300 rounded-lg
                 hover:border-neon-blue hover:text-neon-blue transition-colors"
          :disabled="loadingMore"
          @click="$emit('load-more')"
        >
          <span v-if="loadingMore" class="flex items-center gap-2">
            <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            جاري التحميل...
          </span>
          <span v-else>عرض المزيد</span>
        </button>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import ReviewCard from './ReviewCard.vue'
import RatingsAggregation from './RatingsAggregation.vue'
import type { Review } from '@/lib/types'

const props = withDefaults(defineProps<{
  reviews: Review[]
  totalReviews: number
  avgRating?: number
  loading?: boolean
  loadingMore?: boolean
  hasMore?: boolean
  showAddButton?: boolean
  initialSort?: string
}>(), {
  avgRating: 0,
  loading: false,
  loadingMore: false,
  hasMore: false,
  showAddButton: true,
  initialSort: 'newest'
})

const emit = defineEmits<{
  (e: 'sort-change', sort: string): void
  (e: 'load-more'): void
  (e: 'add-review'): void
  (e: 'helpful', reviewId: number): void
  (e: 'report', reviewId: number): void
}>()

const currentSort = ref(props.initialSort)

watch(() => props.initialSort, (newSort) => {
  currentSort.value = newSort
})

function handleSortChange() {
  emit('sort-change', currentSort.value)
}

function handleHelpful(reviewId: number) {
  emit('helpful', reviewId)
}

function handleReport(reviewId: number) {
  emit('report', reviewId)
}
</script>

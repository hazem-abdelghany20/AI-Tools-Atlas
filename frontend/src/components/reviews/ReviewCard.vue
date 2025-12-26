<template>
  <article class="review-card bg-dark-surface border border-gray-800 rounded-xl p-5" dir="rtl">
    <!-- Header: User Info & Rating -->
    <div class="flex items-start justify-between gap-4 mb-4">
      <!-- User Info -->
      <div class="flex items-center gap-3">
        <!-- Avatar -->
        <div class="w-10 h-10 rounded-full bg-neon-blue/20 flex items-center justify-center flex-shrink-0">
          <span class="text-neon-blue font-bold">
            {{ (review.user?.display_name && review.user.display_name.length > 0) ? review.user.display_name.charAt(0).toUpperCase() : 'U' }}
          </span>
        </div>
        <div>
          <div class="font-medium text-white">
            {{ review.user?.display_name || 'مستخدم' }}
          </div>
          <div class="flex items-center gap-2 text-sm text-gray-400">
            <span v-if="review.reviewer_role">{{ review.reviewer_role }}</span>
            <span v-if="review.reviewer_role && review.company_size" class="text-gray-600">|</span>
            <span v-if="review.company_size">{{ review.company_size }}</span>
          </div>
        </div>
      </div>

      <!-- Rating & Date -->
      <div class="text-left">
        <div class="flex items-center gap-1 mb-1">
          <template v-for="star in 5" :key="star">
            <svg
              class="w-4 h-4"
              :class="star <= review.rating_overall ? 'text-yellow-400' : 'text-gray-600'"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
          </template>
        </div>
        <div class="text-xs text-gray-500">
          {{ formatDate(review.created_at) }}
        </div>
      </div>
    </div>

    <!-- Use Case Badge -->
    <div v-if="review.primary_use_case" class="mb-4">
      <span class="px-3 py-1 text-xs bg-neon-blue/10 text-neon-blue rounded-full">
        {{ review.primary_use_case }}
      </span>
    </div>

    <!-- Pros Section -->
    <div v-if="review.pros" class="mb-4 p-3 rounded-lg bg-green-500/5 border border-green-500/20">
      <div class="flex items-center gap-2 mb-2">
        <svg class="w-4 h-4 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <span class="text-sm font-medium text-green-400">الإيجابيات</span>
      </div>
      <p class="text-gray-300 text-sm">{{ review.pros }}</p>
    </div>

    <!-- Cons Section -->
    <div v-if="review.cons" class="mb-4 p-3 rounded-lg bg-red-500/5 border border-red-500/20">
      <div class="flex items-center gap-2 mb-2">
        <svg class="w-4 h-4 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        <span class="text-sm font-medium text-red-400">السلبيات</span>
      </div>
      <p class="text-gray-300 text-sm">{{ review.cons }}</p>
    </div>

    <!-- Usage Context -->
    <div v-if="review.usage_context" class="mb-4">
      <p class="text-gray-400 text-sm italic">"{{ review.usage_context }}"</p>
    </div>

    <!-- Footer: Helpful & Report Buttons -->
    <div class="flex items-center justify-between pt-3 border-t border-gray-800">
      <button
        class="flex items-center gap-2 text-gray-400 hover:text-neon-blue transition-colors text-sm"
        @click="handleHelpful"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" />
        </svg>
        <span>مفيد؟</span>
        <span v-if="review.helpful_count > 0" class="text-gray-500">({{ review.helpful_count }})</span>
      </button>
      <button
        class="flex items-center gap-2 text-gray-500 hover:text-red-400 transition-colors text-sm"
        @click="emit('report', review.id)"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <span>إبلاغ</span>
      </button>
    </div>
  </article>
</template>

<script setup lang="ts">
import type { Review } from '@/lib/types'

const props = defineProps<{
  review: Review
}>()

const emit = defineEmits<{
  (e: 'helpful', reviewId: number): void
  (e: 'report', reviewId: number): void
}>()

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return ''

  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return 'اليوم'
  if (diffDays === 1) return 'أمس'
  if (diffDays < 7) return `منذ ${diffDays} أيام`
  if (diffDays < 30) return `منذ ${Math.floor(diffDays / 7)} أسابيع`
  if (diffDays < 365) return `منذ ${Math.floor(diffDays / 30)} شهور`
  return `منذ ${Math.floor(diffDays / 365)} سنوات`
}

function handleHelpful() {
  // Future feature - emit for parent to handle when implemented
  emit('helpful', props.review.id)
}
</script>

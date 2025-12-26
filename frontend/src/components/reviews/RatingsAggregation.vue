<template>
  <div class="ratings-aggregation bg-dark-bg rounded-xl p-6 mb-6" dir="rtl">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <!-- Left: Overall Rating -->
      <div class="text-center md:text-right">
        <!-- Large Rating Display -->
        <div class="flex items-center justify-center md:justify-start gap-4 mb-4">
          <div class="text-5xl font-bold text-white">
            {{ avgRating.toFixed(1) }}
          </div>
          <div>
            <!-- Stars -->
            <div class="flex items-center gap-1 mb-1">
              <template v-for="star in 5" :key="star">
                <svg
                  class="w-6 h-6"
                  :class="getStarClass(star, avgRating)"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                </svg>
              </template>
            </div>
            <p class="text-gray-400 text-sm">
              {{ reviewCount }} مراجعة
            </p>
          </div>
        </div>

        <!-- Rating Distribution -->
        <div class="space-y-2 mt-6">
          <div v-for="stars in [5, 4, 3, 2, 1]" :key="stars" class="flex items-center gap-2">
            <span class="text-sm text-gray-400 w-6">{{ stars }}</span>
            <svg class="w-4 h-4 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
            <div class="flex-1 h-2 bg-gray-800 rounded-full overflow-hidden">
              <div
                class="h-full bg-neon-blue rounded-full transition-all duration-500"
                :style="{ width: getDistributionWidth(stars) + '%' }"
              ></div>
            </div>
            <span class="text-xs text-gray-500 w-8 text-left">
              {{ getDistributionCount(stars) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Right: Dimension Ratings -->
      <div v-if="hasDimensionRatings" class="space-y-4">
        <h3 class="text-sm font-medium text-gray-400 mb-4">تقييمات تفصيلية</h3>

        <DimensionBar label="سهولة الاستخدام" :rating="dimensionRatings.easeOfUse" />
        <DimensionBar label="القيمة مقابل السعر" :rating="dimensionRatings.value" />
        <DimensionBar label="الدقة" :rating="dimensionRatings.accuracy" />
        <DimensionBar label="السرعة" :rating="dimensionRatings.speed" />
        <DimensionBar label="الدعم الفني" :rating="dimensionRatings.support" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Review } from '@/lib/types'
import DimensionBar from './DimensionBar.vue'

const props = defineProps<{
  avgRating: number
  reviewCount: number
  reviews: Review[]
}>()

// Calculate rating distribution from reviews
const ratingDistribution = computed(() => {
  const distribution: Record<number, number> = { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 }

  props.reviews.forEach(review => {
    const rating = Math.min(5, Math.max(1, review.rating_overall))
    distribution[rating]++
  })

  return distribution
})

// Calculate dimension ratings averages
const dimensionRatings = computed(() => {
  const sums = {
    easeOfUse: 0,
    value: 0,
    accuracy: 0,
    speed: 0,
    support: 0
  }
  const counts = {
    easeOfUse: 0,
    value: 0,
    accuracy: 0,
    speed: 0,
    support: 0
  }

  props.reviews.forEach(review => {
    if (review.rating_ease_of_use) {
      sums.easeOfUse += review.rating_ease_of_use
      counts.easeOfUse++
    }
    if (review.rating_value) {
      sums.value += review.rating_value
      counts.value++
    }
    if (review.rating_accuracy) {
      sums.accuracy += review.rating_accuracy
      counts.accuracy++
    }
    if (review.rating_speed) {
      sums.speed += review.rating_speed
      counts.speed++
    }
    if (review.rating_support) {
      sums.support += review.rating_support
      counts.support++
    }
  })

  return {
    easeOfUse: counts.easeOfUse > 0 ? sums.easeOfUse / counts.easeOfUse : 0,
    value: counts.value > 0 ? sums.value / counts.value : 0,
    accuracy: counts.accuracy > 0 ? sums.accuracy / counts.accuracy : 0,
    speed: counts.speed > 0 ? sums.speed / counts.speed : 0,
    support: counts.support > 0 ? sums.support / counts.support : 0
  }
})

const hasDimensionRatings = computed(() => {
  const r = dimensionRatings.value
  return r.easeOfUse > 0 || r.value > 0 || r.accuracy > 0 || r.speed > 0 || r.support > 0
})

function getStarClass(star: number, rating: number): string {
  if (star <= Math.floor(rating)) {
    return 'text-yellow-400'
  } else if (star === Math.ceil(rating) && rating % 1 !== 0) {
    return 'text-yellow-400/50'
  }
  return 'text-gray-600'
}

function getDistributionWidth(stars: number): number {
  const total = Object.values(ratingDistribution.value).reduce((a, b) => a + b, 0)
  if (total === 0) return 0
  return (ratingDistribution.value[stars] / total) * 100
}

function getDistributionCount(stars: number): number {
  return ratingDistribution.value[stars] || 0
}
</script>

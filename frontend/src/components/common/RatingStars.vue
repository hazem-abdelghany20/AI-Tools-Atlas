<template>
  <div class="rating-stars flex items-center gap-1" dir="ltr">
    <span class="text-yellow-400 font-semibold">{{ formattedRating }}</span>
    <svg
      v-for="i in 5"
      :key="i"
      class="w-4 h-4"
      :class="getStarClass(i)"
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="currentColor"
    >
      <path
        fill-rule="evenodd"
        d="M10.788 3.21c.448-1.077 1.976-1.077 2.424 0l2.082 5.007 5.404.433c1.164.093 1.636 1.545.749 2.305l-4.117 3.527 1.257 5.273c.271 1.136-.964 2.033-1.96 1.425L12 18.354 7.373 21.18c-.996.608-2.231-.29-1.96-1.425l1.257-5.273-4.117-3.527c-.887-.76-.415-2.212.749-2.305l5.404-.433 2.082-5.006z"
        clip-rule="evenodd"
      />
    </svg>
    <span v-if="showCount && reviewCount > 0" class="text-gray-500 text-sm mr-1">
      ({{ reviewCount }} مراجعة)
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  rating: number
  reviewCount?: number
  showCount?: boolean
}>(), {
  reviewCount: 0,
  showCount: true
})

const formattedRating = computed(() => props.rating.toFixed(1))

function getStarClass(starIndex: number): string {
  const fullStars = Math.floor(props.rating)
  const hasPartialStar = props.rating % 1 >= 0.5

  if (starIndex <= fullStars) {
    return 'text-yellow-400'
  } else if (starIndex === fullStars + 1 && hasPartialStar) {
    return 'text-yellow-400/50'
  } else {
    return 'text-gray-600'
  }
}
</script>

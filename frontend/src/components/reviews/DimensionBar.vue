<template>
  <div class="dimension-bar flex items-center gap-3">
    <span class="text-sm text-gray-400 w-32 flex-shrink-0">{{ label }}</span>
    <div class="flex-1 h-2 bg-gray-800 rounded-full overflow-hidden">
      <div
        class="h-full bg-neon-blue rounded-full transition-all duration-500"
        :style="{ width: (rating / 5) * 100 + '%' }"
      ></div>
    </div>
    <span class="text-sm text-white w-8 text-left">
      {{ rating > 0 ? rating.toFixed(1) : '-' }}
    </span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  rating: number
}>()
</script>

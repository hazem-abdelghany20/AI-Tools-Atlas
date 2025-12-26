<template>
  <div class="spinner-container">
    <div class="spinner" :class="sizeClass"></div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  size?: 'sm' | 'md' | 'lg'
}>()

const sizeClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'w-4 h-4'
    case 'lg':
      return 'w-12 h-12'
    default:
      return 'w-8 h-8'
  }
})
</script>

<style scoped>
.spinner-container {
  @apply flex justify-center items-center;
}

.spinner {
  border: 3px solid rgba(59, 130, 246, 0.3);
  border-top-color: #3B82F6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>

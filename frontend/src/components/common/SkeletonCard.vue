<template>
  <div
    class="skeleton-card animate-pulse bg-dark-surface"
    :style="{ height: height, width: width }"
    :class="{ 'rounded-lg': rounded }"
  ></div>
</template>

<script setup lang="ts">
defineProps<{
  height?: string
  width?: string
  rounded?: boolean
}>()
</script>

<style scoped>
.skeleton-card {
  background: linear-gradient(90deg, #0A0B10 25%, #1F2933 50%, #0A0B10 75%);
  background-size: 200% 100%;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}
</style>

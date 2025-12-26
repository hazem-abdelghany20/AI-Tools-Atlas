<template>
  <div class="dimension-rating flex items-center justify-between gap-2">
    <span class="text-sm text-gray-400 flex-shrink-0">{{ label }}</span>
    <div class="flex items-center gap-0.5">
      <button
        v-for="star in 5"
        :key="star"
        type="button"
        @click="$emit('update:modelValue', star)"
        class="p-0.5 transition-colors"
      >
        <svg
          class="w-4 h-4"
          :class="star <= modelValue ? 'text-yellow-400' : 'text-gray-600 hover:text-gray-500'"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  modelValue: number
}>()

defineEmits<{
  (e: 'update:modelValue', value: number): void
}>()
</script>

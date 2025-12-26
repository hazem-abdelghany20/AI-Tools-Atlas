<template>
  <div class="compare-bar fixed bottom-0 left-0 right-0 bg-dark-surface border-t border-gray-800 p-4 z-40" dir="rtl">
    <div class="max-w-5xl mx-auto flex items-center justify-between">
      <!-- Selected Tools -->
      <div class="flex items-center gap-3">
        <span class="text-gray-400 text-sm">للمقارنة ({{ tools.length }}/4):</span>
        <div class="flex items-center gap-2">
          <div
            v-for="tool in tools"
            :key="tool.id"
            class="flex items-center gap-2 bg-dark-bg px-3 py-1.5 rounded-full"
          >
            <img
              v-if="tool.logo_url"
              :src="tool.logo_url"
              :alt="tool.name"
              class="w-5 h-5 rounded"
            />
            <span class="text-white text-sm">{{ tool.name }}</span>
            <button
              @click="$emit('remove', tool)"
              class="text-gray-500 hover:text-red-400 transition-colors"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-3">
        <button
          @click="$emit('clear')"
          class="text-gray-400 hover:text-white text-sm transition-colors"
        >
          مسح الكل
        </button>
        <button
          @click="$emit('compare')"
          :disabled="tools.length < 2"
          class="px-6 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg
                 hover:bg-neon-blue-hover transition-colors disabled:opacity-50"
        >
          مقارنة ({{ tools.length }})
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Tool } from '@/lib/types'

defineProps<{
  tools: Tool[]
}>()

defineEmits<{
  (e: 'remove', tool: Tool): void
  (e: 'clear'): void
  (e: 'compare'): void
}>()
</script>

<template>
  <section v-if="hasFeatures" class="tool-features bg-dark-surface rounded-lg p-6 md:p-8">
    <h2 class="text-xl font-bold mb-6">المميزات</h2>

    <!-- Array of features -->
    <div v-if="featuresList.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="(feature, index) in featuresList"
        :key="index"
        class="feature-card p-4 bg-dark-border rounded-lg"
      >
        <div class="flex items-start gap-3">
          <div class="flex-shrink-0 w-10 h-10 rounded-lg bg-neon-blue/20 flex items-center justify-center">
            <svg v-if="!feature.icon" class="w-5 h-5 text-neon-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            <span v-else class="text-lg">{{ feature.icon }}</span>
          </div>
          <div>
            <h3 class="font-semibold text-white mb-1">{{ feature.name }}</h3>
            <p v-if="feature.description" class="text-sm text-gray-400">{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Text features (comma-separated or line-separated) -->
    <div v-else-if="featuresText" class="space-y-3">
      <div
        v-for="(feature, index) in featuresTextList"
        :key="index"
        class="flex items-start gap-3"
      >
        <svg class="w-5 h-5 text-neon-blue flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <span class="text-gray-300">{{ feature }}</span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Tool, ToolFeature } from '../../lib/types'

const props = defineProps<{
  tool: Tool
}>()

const hasFeatures = computed(() => {
  if (!props.tool.features) return false
  if (Array.isArray(props.tool.features)) return props.tool.features.length > 0
  return props.tool.features.trim().length > 0
})

const featuresList = computed((): ToolFeature[] => {
  if (!props.tool.features || typeof props.tool.features === 'string') return []
  return props.tool.features
})

const featuresText = computed(() => {
  if (!props.tool.features || typeof props.tool.features !== 'string') return ''
  return props.tool.features
})

const featuresTextList = computed(() => {
  if (!featuresText.value) return []
  // Support both newline and comma-separated features
  const separator = featuresText.value.includes('\n') ? '\n' : ','
  return featuresText.value.split(separator).map(f => f.trim()).filter(Boolean)
})
</script>

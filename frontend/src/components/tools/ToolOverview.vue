<template>
  <section class="tool-overview bg-dark-surface rounded-lg p-6 md:p-8 space-y-6">
    <!-- Best For -->
    <div v-if="tool.best_for" class="best-for">
      <h2 class="text-lg font-bold text-neon-blue mb-2">الأفضل لـ</h2>
      <p class="text-white text-lg leading-relaxed">{{ tool.best_for }}</p>
    </div>

    <!-- Description -->
    <div v-if="tool.description" class="description">
      <h2 class="text-lg font-bold text-gray-300 mb-3">نظرة عامة</h2>
      <div class="text-gray-300 leading-relaxed space-y-4">
        <p v-for="(paragraph, index) in descriptionParagraphs" :key="index">
          {{ paragraph }}
        </p>
      </div>
    </div>

    <!-- Primary Use Cases -->
    <div v-if="tool.primary_use_cases" class="use-cases">
      <h2 class="text-lg font-bold text-gray-300 mb-3">حالات الاستخدام الرئيسية</h2>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="(useCase, index) in useCasesList"
          :key="index"
          class="px-3 py-2 bg-dark-border rounded-lg text-gray-200"
        >
          {{ useCase }}
        </span>
      </div>
    </div>

    <!-- Target Roles -->
    <div v-if="tool.target_roles" class="target-roles">
      <h2 class="text-lg font-bold text-gray-300 mb-3">مناسب لـ</h2>
      <div class="flex flex-wrap gap-2">
        <span
          v-for="(role, index) in targetRolesList"
          :key="index"
          class="px-3 py-2 bg-primary-500/20 text-primary-500 rounded-lg"
        >
          {{ role }}
        </span>
      </div>
    </div>

    <!-- Platforms -->
    <div v-if="tool.platforms" class="platforms">
      <h2 class="text-lg font-bold text-gray-300 mb-3">المنصات المدعومة</h2>
      <div class="flex flex-wrap gap-3">
        <div
          v-for="(platform, index) in platformsList"
          :key="index"
          class="flex items-center gap-2 px-4 py-2 bg-dark-border rounded-lg"
        >
          <component :is="getPlatformIcon(platform)" class="w-5 h-5 text-gray-400" />
          <span class="text-gray-200">{{ platform }}</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, h } from 'vue'
import type { Tool } from '../../lib/types'

const props = defineProps<{
  tool: Tool
}>()

const descriptionParagraphs = computed(() => {
  if (!props.tool.description) return []
  return props.tool.description.split('\n').filter(p => p.trim())
})

const useCasesList = computed(() => {
  if (!props.tool.primary_use_cases) return []
  return props.tool.primary_use_cases.split(',').map(s => s.trim()).filter(Boolean)
})

const targetRolesList = computed(() => {
  if (!props.tool.target_roles) return []
  return props.tool.target_roles.split(',').map(s => s.trim()).filter(Boolean)
})

const platformsList = computed(() => {
  if (!props.tool.platforms) return []
  return props.tool.platforms.split(',').map(s => s.trim()).filter(Boolean)
})

// Platform icons as simple SVG components
const WebIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9' })
])

const MobileIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z' })
])

const DesktopIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z' })
])

const ApiIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4' })
])

const DefaultIcon = () => h('svg', { fill: 'none', viewBox: '0 0 24 24', stroke: 'currentColor' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 6h16M4 12h16M4 18h16' })
])

const getPlatformIcon = (platform: string) => {
  const normalized = platform.toLowerCase()
  if (normalized.includes('web')) return WebIcon
  if (normalized.includes('mobile') || normalized.includes('ios') || normalized.includes('android')) return MobileIcon
  if (normalized.includes('desktop') || normalized.includes('windows') || normalized.includes('mac')) return DesktopIcon
  if (normalized.includes('api')) return ApiIcon
  return DefaultIcon
}
</script>

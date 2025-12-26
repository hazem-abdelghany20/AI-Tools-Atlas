<template>
  <section class="tool-pricing bg-dark-surface rounded-lg p-6 md:p-8">
    <h2 class="text-xl font-bold mb-6">الأسعار</h2>

    <div class="space-y-4">
      <!-- Free Tier Badge -->
      <div v-if="tool.has_free_tier" class="flex items-center gap-2 text-green-400">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <span class="font-medium">يتوفر خطة مجانية</span>
      </div>

      <!-- Pricing Summary -->
      <div v-if="tool.pricing_summary" class="pricing-summary">
        <div class="bg-dark-border rounded-lg p-4">
          <div class="flex items-center gap-2 mb-2">
            <svg class="w-5 h-5 text-neon-blue" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-gray-400 text-sm">معلومات التسعير</span>
          </div>
          <p class="text-white text-lg">{{ tool.pricing_summary }}</p>
        </div>
      </div>

      <!-- No Pricing Info -->
      <div v-else-if="!tool.has_free_tier" class="text-gray-400">
        <p>الرجاء زيارة الموقع الرسمي للأسعار</p>
      </div>

      <!-- Link to Official Pricing -->
      <div v-if="tool.official_url" class="mt-4">
        <a
          :href="pricingUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 text-neon-blue hover:text-neon-blue-hover transition-colors"
        >
          <span>عرض تفاصيل الأسعار</span>
          <svg class="w-4 h-4 rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
          </svg>
        </a>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Tool } from '../../lib/types'

const props = defineProps<{
  tool: Tool
}>()

const pricingUrl = computed(() => {
  if (!props.tool.official_url) return ''
  // Try to construct pricing page URL, fallback to official URL
  const baseUrl = props.tool.official_url.replace(/\/$/, '')
  return `${baseUrl}/pricing`
})
</script>

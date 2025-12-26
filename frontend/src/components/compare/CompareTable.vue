<template>
  <div class="compare-table overflow-x-auto" dir="rtl">
    <table class="w-full border-collapse">
      <!-- Header Row with Tool Info -->
      <thead>
        <tr class="border-b border-gray-800">
          <th class="p-4 text-right text-gray-400 font-medium w-40 sticky right-0 bg-dark-bg z-10">
            المقارنة
          </th>
          <th
            v-for="tool in tools"
            :key="tool.id"
            class="p-4 text-center min-w-[200px]"
          >
            <div class="flex flex-col items-center gap-3">
              <!-- Remove Button -->
              <button
                @click="$emit('remove', tool)"
                class="absolute top-2 left-2 p-1 text-gray-500 hover:text-red-400 transition-colors"
                title="إزالة من المقارنة"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>

              <!-- Tool Logo -->
              <img
                v-if="tool.logo_url"
                :src="tool.logo_url"
                :alt="tool.name"
                class="w-16 h-16 rounded-xl object-cover"
              />
              <div
                v-else
                class="w-16 h-16 rounded-xl bg-neon-blue/20 flex items-center justify-center"
              >
                <span class="text-neon-blue font-bold text-2xl">{{ tool.name.charAt(0) }}</span>
              </div>

              <!-- Tool Name -->
              <h3 class="font-semibold text-white text-lg">{{ tool.name }}</h3>

              <!-- Tagline -->
              <p class="text-gray-400 text-sm line-clamp-2">{{ tool.tagline }}</p>
            </div>
          </th>
        </tr>
      </thead>

      <tbody>
        <!-- Best For Row -->
        <tr class="border-b border-gray-800/50">
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-bg">الأفضل لـ</td>
          <td v-for="tool in tools" :key="'best-' + tool.id" class="p-4 text-center text-gray-300">
            {{ tool.best_for || '-' }}
          </td>
        </tr>

        <!-- Category Row -->
        <tr class="border-b border-gray-800/50 bg-dark-surface/30">
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-surface/30">التصنيف</td>
          <td v-for="tool in tools" :key="'cat-' + tool.id" class="p-4 text-center">
            <span v-if="tool.primary_category" class="px-3 py-1 bg-gray-800 text-gray-300 rounded-full text-sm">
              {{ tool.primary_category.name }}
            </span>
            <span v-else class="text-gray-500">-</span>
          </td>
        </tr>

        <!-- Rating Row -->
        <tr class="border-b border-gray-800/50">
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-bg">التقييم</td>
          <td v-for="tool in tools" :key="'rating-' + tool.id" class="p-4 text-center">
            <div class="flex flex-col items-center gap-1">
              <div class="flex items-center gap-1">
                <span class="text-yellow-400 text-lg font-bold">{{ tool.avg_rating_overall?.toFixed(1) || '0.0' }}</span>
                <svg class="w-5 h-5 text-yellow-400" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                </svg>
              </div>
              <span class="text-gray-500 text-sm">({{ tool.review_count || 0 }} تقييم)</span>
            </div>
          </td>
        </tr>

        <!-- Pricing Row -->
        <tr class="border-b border-gray-800/50 bg-dark-surface/30">
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-surface/30">التسعير</td>
          <td v-for="tool in tools" :key="'pricing-' + tool.id" class="p-4 text-center">
            <div class="flex flex-col items-center gap-2">
              <span
                class="px-3 py-1 rounded-full text-sm font-medium"
                :class="tool.has_free_tier ? 'bg-green-500/20 text-green-400' : 'bg-gray-700 text-gray-300'"
              >
                {{ tool.has_free_tier ? 'يوجد خطة مجانية' : 'مدفوع' }}
              </span>
              <span v-if="tool.pricing_summary" class="text-gray-400 text-sm">
                {{ tool.pricing_summary }}
              </span>
            </div>
          </td>
        </tr>

        <!-- Platforms Row -->
        <tr class="border-b border-gray-800/50">
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-bg">المنصات</td>
          <td v-for="tool in tools" :key="'platforms-' + tool.id" class="p-4 text-center text-gray-300">
            {{ tool.platforms || '-' }}
          </td>
        </tr>

        <!-- Target Roles Row -->
        <tr class="border-b border-gray-800/50 bg-dark-surface/30">
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-surface/30">الفئة المستهدفة</td>
          <td v-for="tool in tools" :key="'roles-' + tool.id" class="p-4 text-center text-gray-300">
            {{ tool.target_roles || '-' }}
          </td>
        </tr>

        <!-- Bookmark Count Row -->
        <tr class="border-b border-gray-800/50">
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-bg">عدد الحفظ</td>
          <td v-for="tool in tools" :key="'bookmarks-' + tool.id" class="p-4 text-center">
            <div class="flex items-center justify-center gap-2 text-gray-300">
              <svg class="w-5 h-5 text-red-400" fill="currentColor" viewBox="0 0 24 24">
                <path d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
              </svg>
              <span>{{ tool.bookmark_count || 0 }}</span>
            </div>
          </td>
        </tr>

        <!-- Badges Row -->
        <tr class="border-b border-gray-800/50 bg-dark-surface/30">
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-surface/30">الشارات</td>
          <td v-for="tool in tools" :key="'badges-' + tool.id" class="p-4 text-center">
            <div v-if="tool.badges && tool.badges.length > 0" class="flex flex-wrap justify-center gap-2">
              <span
                v-for="badge in tool.badges"
                :key="badge.id"
                class="px-2 py-1 bg-neon-blue/20 text-neon-blue rounded-full text-xs"
              >
                {{ badge.name }}
              </span>
            </div>
            <span v-else class="text-gray-500">-</span>
          </td>
        </tr>

        <!-- Actions Row -->
        <tr>
          <td class="p-4 text-gray-400 font-medium sticky right-0 bg-dark-bg">الإجراءات</td>
          <td v-for="tool in tools" :key="'actions-' + tool.id" class="p-4 text-center">
            <div class="flex flex-col items-center gap-2">
              <router-link
                :to="`/tools/${tool.slug}`"
                class="px-4 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg hover:bg-neon-blue-hover transition-colors text-sm"
              >
                عرض التفاصيل
              </router-link>
              <a
                v-if="tool.official_url"
                :href="tool.official_url"
                target="_blank"
                rel="noopener noreferrer"
                class="px-4 py-2 border border-gray-700 text-gray-300 rounded-lg hover:border-neon-blue hover:text-neon-blue transition-colors text-sm"
              >
                زيارة الموقع
              </a>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { Tool } from '@/lib/types'

defineProps<{
  tools: Tool[]
}>()

defineEmits<{
  (e: 'remove', tool: Tool): void
}>()
</script>

<style scoped>
.compare-table {
  min-width: 100%;
}

.compare-table th,
.compare-table td {
  position: relative;
}

.compare-table th:first-child,
.compare-table td:first-child {
  min-width: 160px;
}
</style>

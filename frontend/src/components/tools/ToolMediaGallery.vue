<template>
  <section v-if="hasMedia" class="tool-media-gallery bg-dark-surface rounded-lg p-6 md:p-8">
    <h2 class="text-xl font-bold mb-6">صور وفيديوهات</h2>

    <!-- Screenshots Grid -->
    <div v-if="screenshots.length > 0" class="mb-8">
      <h3 v-if="videos.length > 0" class="text-gray-400 text-sm mb-4">لقطات الشاشة</h3>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="(screenshot, index) in screenshots"
          :key="screenshot.id"
          class="screenshot-item relative cursor-pointer group overflow-hidden rounded-lg bg-dark-border aspect-video"
          @click="openLightbox(index)"
        >
          <img
            :src="screenshot.thumbnail_url || screenshot.url"
            :alt="`لقطة شاشة ${index + 1}`"
            class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            loading="lazy"
          />
          <div class="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-colors flex items-center justify-center">
            <svg class="w-10 h-10 text-white opacity-0 group-hover:opacity-100 transition-opacity" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Videos Grid -->
    <div v-if="videos.length > 0">
      <h3 v-if="screenshots.length > 0" class="text-gray-400 text-sm mb-4">مقاطع الفيديو</h3>
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <YouTubeEmbed
          v-for="video in videos"
          :key="video.id"
          :url="video.url"
          :alt="video.thumbnail_url ? `فيديو توضيحي` : undefined"
        />
      </div>
    </div>

    <!-- Lightbox -->
    <Lightbox
      ref="lightboxRef"
      :images="screenshotUrls"
      alt="لقطة شاشة"
    />
  </section>

  <!-- Empty State (optional - can be hidden instead) -->
  <section v-else-if="showEmptyState" class="tool-media-gallery bg-dark-surface rounded-lg p-6 md:p-8">
    <h2 class="text-xl font-bold mb-4">صور وفيديوهات</h2>
    <p class="text-gray-400">لا توجد صور أو فيديوهات متاحة</p>
  </section>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Tool } from '../../lib/types'
import Lightbox from '../common/Lightbox.vue'
import YouTubeEmbed from '../common/YouTubeEmbed.vue'

const props = withDefaults(defineProps<{
  tool: Tool
  showEmptyState?: boolean
}>(), {
  showEmptyState: false
})

const lightboxRef = ref<InstanceType<typeof Lightbox> | null>(null)

const hasMedia = computed(() => {
  return props.tool.media && props.tool.media.length > 0
})

const sortedMedia = computed(() => {
  if (!props.tool.media) return []
  return [...props.tool.media].sort((a, b) => a.display_order - b.display_order)
})

const screenshots = computed(() => {
  return sortedMedia.value.filter(m => m.type === 'screenshot')
})

const videos = computed(() => {
  return sortedMedia.value.filter(m => m.type === 'video')
})

const screenshotUrls = computed(() => {
  return screenshots.value.map(s => s.url)
})

const openLightbox = (index: number) => {
  lightboxRef.value?.open(index)
}
</script>

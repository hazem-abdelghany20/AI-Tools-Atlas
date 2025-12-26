<template>
  <div class="youtube-embed relative bg-dark-border rounded-lg overflow-hidden" :style="{ aspectRatio: '16/9' }">
    <!-- Thumbnail with play button (lazy load) -->
    <div v-if="!loaded" class="relative w-full h-full cursor-pointer" @click="loadVideo">
      <img
        v-if="thumbnailUrl"
        :src="thumbnailUrl"
        :alt="alt"
        class="w-full h-full object-cover"
        loading="lazy"
      />
      <div class="absolute inset-0 flex items-center justify-center bg-black/30">
        <div class="play-button w-16 h-16 bg-red-600 rounded-full flex items-center justify-center hover:bg-red-700 transition-colors">
          <svg class="w-8 h-8 text-white ms-1" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z" />
          </svg>
        </div>
      </div>
    </div>

    <!-- Actual iframe (loads on click) -->
    <iframe
      v-else
      :src="embedUrl"
      class="w-full h-full"
      frameborder="0"
      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
      allowfullscreen
    ></iframe>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  url: string
  alt?: string
}>()

const loaded = ref(false)

const videoId = computed(() => {
  // Extract video ID from various YouTube URL formats
  const patterns = [
    /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([a-zA-Z0-9_-]{11})/,
    /youtube\.com\/v\/([a-zA-Z0-9_-]{11})/
  ]

  for (const pattern of patterns) {
    const match = props.url.match(pattern)
    if (match) return match[1]
  }

  return null
})

const thumbnailUrl = computed(() => {
  if (!videoId.value) return ''
  return `https://img.youtube.com/vi/${videoId.value}/maxresdefault.jpg`
})

const embedUrl = computed(() => {
  if (!videoId.value) return ''
  return `https://www.youtube.com/embed/${videoId.value}?autoplay=1`
})

const loadVideo = () => {
  loaded.value = true
}
</script>

<style scoped>
.play-button {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}
</style>

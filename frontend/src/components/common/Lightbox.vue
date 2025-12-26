<template>
  <Teleport to="body">
    <Transition name="lightbox">
      <div
        v-if="isOpen"
        class="lightbox fixed inset-0 z-50 flex items-center justify-center"
        @click.self="close"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/90"></div>

        <!-- Content -->
        <div class="relative z-10 max-w-[90vw] max-h-[90vh]">
          <!-- Close button -->
          <button
            @click="close"
            class="absolute -top-12 start-0 text-white hover:text-gray-300 transition-colors"
          >
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>

          <!-- Previous button -->
          <button
            v-if="hasPrevious"
            @click="previous"
            class="nav-button absolute start-0 top-1/2 -translate-y-1/2 -ms-16 text-white hover:text-neon-blue transition-colors"
          >
            <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

          <!-- Image -->
          <img
            :src="currentImage"
            :alt="alt"
            class="max-w-full max-h-[85vh] rounded-lg"
          />

          <!-- Next button -->
          <button
            v-if="hasNext"
            @click="next"
            class="nav-button absolute end-0 top-1/2 -translate-y-1/2 -me-16 text-white hover:text-neon-blue transition-colors"
          >
            <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>

          <!-- Counter -->
          <div v-if="images.length > 1" class="absolute -bottom-10 left-1/2 -translate-x-1/2 text-white text-sm">
            {{ currentIndex + 1 }} / {{ images.length }}
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  images: string[]
  initialIndex?: number
  alt?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const isOpen = ref(false)
const currentIndex = ref(props.initialIndex || 0)

const currentImage = computed(() => props.images[currentIndex.value] || '')
const hasPrevious = computed(() => currentIndex.value > 0)
const hasNext = computed(() => currentIndex.value < props.images.length - 1)

const open = (index: number = 0) => {
  currentIndex.value = index
  isOpen.value = true
  document.body.style.overflow = 'hidden'
}

const close = () => {
  isOpen.value = false
  document.body.style.overflow = ''
  emit('close')
}

const previous = () => {
  if (hasPrevious.value) {
    currentIndex.value--
  }
}

const next = () => {
  if (hasNext.value) {
    currentIndex.value++
  }
}

const handleKeydown = (e: KeyboardEvent) => {
  if (!isOpen.value) return

  switch (e.key) {
    case 'Escape':
      close()
      break
    case 'ArrowLeft':
      next() // RTL: left goes to next
      break
    case 'ArrowRight':
      previous() // RTL: right goes to previous
      break
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
})

// Expose methods for parent component
defineExpose({ open, close })
</script>

<style scoped>
.lightbox-enter-active,
.lightbox-leave-active {
  transition: opacity 0.3s ease;
}

.lightbox-enter-from,
.lightbox-leave-to {
  opacity: 0;
}
</style>

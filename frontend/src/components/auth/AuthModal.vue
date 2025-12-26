<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="$emit('close')"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm"></div>

        <!-- Modal -->
        <div
          class="relative w-full max-w-md bg-dark-surface border border-gray-800 rounded-2xl shadow-xl"
          dir="rtl"
        >
          <!-- Header -->
          <div class="flex items-center justify-between p-6 border-b border-gray-800">
            <h2 class="text-xl font-semibold text-white">
              {{ mode === 'signin' ? 'تسجيل الدخول' : 'إنشاء حساب' }}
            </h2>
            <button
              @click="$emit('close')"
              class="p-2 text-gray-400 hover:text-white transition-colors rounded-lg hover:bg-gray-800"
            >
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Body -->
          <div class="p-6">
            <SignInForm
              v-if="mode === 'signin'"
              @switch="mode = 'signup'"
              @success="handleSuccess"
            />
            <SignUpForm
              v-else
              @switch="mode = 'signin'"
              @success="handleSuccess"
            />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- Success Toast -->
  <Teleport to="body">
    <Transition name="toast">
      <div
        v-if="showSuccessToast"
        class="fixed bottom-6 left-1/2 -translate-x-1/2 px-6 py-3 bg-green-500 text-white rounded-lg shadow-lg z-50"
      >
        {{ successMessage }}
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import SignInForm from './SignInForm.vue'
import SignUpForm from './SignUpForm.vue'

const props = defineProps<{
  isOpen: boolean
  initialMode?: 'signin' | 'signup'
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const mode = ref<'signin' | 'signup'>(props.initialMode || 'signin')
const showSuccessToast = ref(false)
const successMessage = ref('')

// Reset mode when modal opens
watch(() => props.isOpen, (isOpen) => {
  if (isOpen && props.initialMode) {
    mode.value = props.initialMode
  }
})

function handleSuccess() {
  const message = mode.value === 'signin'
    ? 'تم تسجيل الدخول بنجاح'
    : 'تم إنشاء الحساب بنجاح'

  successMessage.value = message
  showSuccessToast.value = true

  emit('close')

  setTimeout(() => {
    showSuccessToast.value = false
  }, 3000)
}

// Handle Escape key to close modal
function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.isOpen) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
/* Modal Transition */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active .relative,
.modal-leave-active .relative {
  transition: transform 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.95);
}

/* Toast Transition */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 20px);
}
</style>

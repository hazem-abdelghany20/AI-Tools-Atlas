<template>
  <div class="fixed inset-0 z-50 overflow-y-auto" dir="rtl">
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black/70" @click="$emit('close')"></div>

    <!-- Modal -->
    <div class="relative min-h-screen flex items-center justify-center p-4">
      <div class="relative bg-dark-surface border border-gray-700 rounded-xl w-full max-w-md">
        <!-- Header -->
        <div class="border-b border-gray-800 px-6 py-4 flex justify-between items-center">
          <h2 class="text-lg font-bold text-white">
            <span class="text-red-400 ml-2">
              <svg class="w-5 h-5 inline-block" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </span>
            الإبلاغ عن محتوى
          </h2>
          <button @click="$emit('close')" class="text-gray-400 hover:text-white">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
          <!-- Reason -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              سبب الإبلاغ <span class="text-red-400">*</span>
            </label>
            <select
              v-model="form.reason"
              required
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white focus:outline-none focus:border-neon-blue"
            >
              <option value="" disabled>اختر السبب</option>
              <option value="spam">محتوى مزعج (Spam)</option>
              <option value="abuse">محتوى مسيء (Abuse)</option>
              <option value="misinformation">معلومات مضللة (Misinformation)</option>
              <option value="other">سبب آخر (Other)</option>
            </select>
          </div>

          <!-- Comment -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              تفاصيل إضافية (اختياري)
            </label>
            <textarea
              v-model="form.comment"
              rows="3"
              maxlength="500"
              class="w-full px-4 py-2 bg-dark-bg border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-neon-blue resize-none"
              placeholder="أضف تفاصيل إضافية حول البلاغ..."
            ></textarea>
            <div class="text-xs text-gray-500 mt-1 text-left">{{ form.comment.length }}/500</div>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="p-3 bg-red-500/20 border border-red-500/50 rounded-lg text-red-400 text-sm">
            {{ error }}
          </div>

          <!-- Success Message -->
          <div v-if="success" class="p-3 bg-green-500/20 border border-green-500/50 rounded-lg text-green-400 text-sm">
            شكراً للإبلاغ. سنراجعه قريباً
          </div>

          <!-- Actions -->
          <div class="flex justify-end gap-3 pt-4 border-t border-gray-800">
            <button
              type="button"
              @click="$emit('close')"
              class="px-4 py-2 text-gray-400 hover:text-white transition-colors"
            >
              إلغاء
            </button>
            <button
              type="submit"
              :disabled="submitting || success"
              class="px-6 py-2 bg-red-500 text-white font-medium rounded-lg hover:bg-red-600 transition-colors disabled:opacity-50 flex items-center gap-2"
            >
              <span v-if="submitting">جاري الإرسال...</span>
              <span v-else>إرسال البلاغ</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { apiClient } from '@/lib/apiClient'

const props = defineProps<{
  reportableType: 'tool' | 'review'
  reportableId: string | number
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'reported'): void
}>()

// Handle Escape key to close modal
function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})

const form = ref({
  reason: '',
  comment: '',
})

const submitting = ref(false)
const error = ref('')
const success = ref(false)

async function handleSubmit() {
  error.value = ''
  submitting.value = true

  try {
    const endpoint = props.reportableType === 'tool'
      ? `/tools/${props.reportableId}/report`
      : `/reviews/${props.reportableId}/report`

    await apiClient.post(endpoint, {
      reason: form.value.reason,
      comment: form.value.comment || undefined,
    })

    success.value = true
    emit('reported')

    // Auto-close after 2 seconds
    setTimeout(() => {
      emit('close')
    }, 2000)
  } catch (err: any) {
    if (err.message?.includes('ALREADY_REPORTED')) {
      error.value = 'لقد أبلغت عن هذا المحتوى مسبقاً اليوم'
    } else {
      error.value = err.message || 'فشل إرسال البلاغ. حاول مرة أخرى'
    }
  } finally {
    submitting.value = false
  }
}

// Reset form when modal opens
watch(() => props.reportableId, () => {
  form.value = { reason: '', comment: '' }
  error.value = ''
  success.value = false
})
</script>

<template>
  <div class="review-form-overlay fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4" @click.self="handleClose">
    <div class="review-form bg-dark-surface rounded-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto" dir="rtl">
      <!-- Header -->
      <div class="flex items-center justify-between p-5 border-b border-gray-800">
        <h2 class="text-xl font-bold text-white">إضافة مراجعة</h2>
        <button
          @click="handleClose"
          class="text-gray-400 hover:text-white transition-colors p-1"
        >
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="p-5 space-y-6">
        <!-- Overall Rating (Required) -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-white">
            التقييم العام <span class="text-red-400">*</span>
          </label>
          <div class="flex items-center gap-1">
            <button
              v-for="star in 5"
              :key="star"
              type="button"
              @click="form.rating_overall = star"
              @mouseenter="hoverRating = star"
              @mouseleave="hoverRating = 0"
              class="p-1 transition-transform hover:scale-110"
            >
              <svg
                class="w-8 h-8"
                :class="star <= (hoverRating || form.rating_overall) ? 'text-yellow-400' : 'text-gray-600'"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
              </svg>
            </button>
          </div>
          <p v-if="errors.rating_overall" class="text-sm text-red-400">{{ errors.rating_overall }}</p>
        </div>

        <!-- Dimension Ratings (Optional) -->
        <div class="space-y-4 p-4 bg-dark-bg rounded-lg">
          <h3 class="text-sm font-medium text-gray-400 mb-3">تقييمات تفصيلية (اختياري)</h3>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <DimensionRating v-model="form.rating_ease_of_use" label="سهولة الاستخدام" />
            <DimensionRating v-model="form.rating_value" label="القيمة مقابل السعر" />
            <DimensionRating v-model="form.rating_accuracy" label="الدقة" />
            <DimensionRating v-model="form.rating_speed" label="السرعة" />
            <DimensionRating v-model="form.rating_support" label="الدعم الفني" />
          </div>
        </div>

        <!-- Pros (Required) -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-white">
            الإيجابيات <span class="text-red-400">*</span>
          </label>
          <textarea
            v-model="form.pros"
            rows="3"
            placeholder="ما الذي أعجبك في هذه الأداة؟"
            class="w-full bg-dark-bg border border-gray-700 rounded-lg px-4 py-3 text-white
                   placeholder-gray-500 focus:border-neon-blue focus:outline-none resize-none"
          ></textarea>
          <div class="flex justify-between text-xs">
            <p v-if="errors.pros" class="text-red-400">{{ errors.pros }}</p>
            <span class="text-gray-500 mr-auto" :class="{ 'text-red-400': charCount(form.pros) > 500 }">{{ charCount(form.pros) }}/500</span>
          </div>
        </div>

        <!-- Cons (Required) -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-white">
            السلبيات <span class="text-red-400">*</span>
          </label>
          <textarea
            v-model="form.cons"
            rows="3"
            placeholder="ما الذي يمكن تحسينه؟"
            class="w-full bg-dark-bg border border-gray-700 rounded-lg px-4 py-3 text-white
                   placeholder-gray-500 focus:border-neon-blue focus:outline-none resize-none"
          ></textarea>
          <div class="flex justify-between text-xs">
            <p v-if="errors.cons" class="text-red-400">{{ errors.cons }}</p>
            <span class="text-gray-500 mr-auto" :class="{ 'text-red-400': charCount(form.cons) > 500 }">{{ charCount(form.cons) }}/500</span>
          </div>
        </div>

        <!-- Dropdowns Row -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- Primary Use Case -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-400">حالة الاستخدام الرئيسية</label>
            <select
              v-model="form.primary_use_case"
              class="w-full bg-dark-bg border border-gray-700 rounded-lg px-4 py-2.5 text-white
                     focus:border-neon-blue focus:outline-none appearance-none"
            >
              <option value="">اختر...</option>
              <option value="content_creation">إنشاء المحتوى</option>
              <option value="coding">البرمجة</option>
              <option value="research">البحث</option>
              <option value="design">التصميم</option>
              <option value="marketing">التسويق</option>
              <option value="data_analysis">تحليل البيانات</option>
              <option value="customer_service">خدمة العملاء</option>
              <option value="other">أخرى</option>
            </select>
          </div>

          <!-- Reviewer Role -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-400">دورك الوظيفي</label>
            <select
              v-model="form.reviewer_role"
              class="w-full bg-dark-bg border border-gray-700 rounded-lg px-4 py-2.5 text-white
                     focus:border-neon-blue focus:outline-none appearance-none"
            >
              <option value="">اختر...</option>
              <option value="developer">مطور</option>
              <option value="designer">مصمم</option>
              <option value="manager">مدير</option>
              <option value="marketer">مسوق</option>
              <option value="writer">كاتب محتوى</option>
              <option value="data_analyst">محلل بيانات</option>
              <option value="student">طالب</option>
              <option value="freelancer">مستقل</option>
              <option value="other">أخرى</option>
            </select>
          </div>

          <!-- Company Size -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-400">حجم الشركة</label>
            <select
              v-model="form.company_size"
              class="w-full bg-dark-bg border border-gray-700 rounded-lg px-4 py-2.5 text-white
                     focus:border-neon-blue focus:outline-none appearance-none"
            >
              <option value="">اختر...</option>
              <option value="1-10">1-10 موظفين</option>
              <option value="11-50">11-50 موظف</option>
              <option value="51-200">51-200 موظف</option>
              <option value="201-500">201-500 موظف</option>
              <option value="500+">أكثر من 500</option>
              <option value="individual">فردي</option>
            </select>
          </div>
        </div>

        <!-- Usage Context -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-400">سياق الاستخدام</label>
          <textarea
            v-model="form.usage_context"
            rows="2"
            placeholder="كيف تستخدم هذه الأداة؟ (اختياري)"
            class="w-full bg-dark-bg border border-gray-700 rounded-lg px-4 py-3 text-white
                   placeholder-gray-500 focus:border-neon-blue focus:outline-none resize-none"
          ></textarea>
        </div>

        <!-- Error Message -->
        <div v-if="submitError" class="p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
          <p class="text-red-400 text-sm">{{ submitError }}</p>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3 pt-4 border-t border-gray-800">
          <button
            type="button"
            @click="handleClose"
            class="px-6 py-2.5 text-gray-400 hover:text-white transition-colors"
          >
            إلغاء
          </button>
          <button
            type="submit"
            :disabled="submitting"
            class="px-6 py-2.5 bg-neon-blue text-dark-bg font-medium rounded-lg
                   hover:bg-neon-blue-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed
                   flex items-center gap-2"
          >
            <svg v-if="submitting" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ submitting ? 'جاري النشر...' : 'نشر المراجعة' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { apiClient } from '@/lib/apiClient'
import DimensionRating from './DimensionRating.vue'

const props = defineProps<{
  toolSlug: string
}>()

// Count Unicode characters, not bytes
function charCount(str: string): number {
  return [...str].length
}

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'success'): void
}>()

const hoverRating = ref(0)
const submitting = ref(false)
const submitError = ref('')

const form = reactive({
  rating_overall: 0,
  rating_ease_of_use: 0,
  rating_value: 0,
  rating_accuracy: 0,
  rating_speed: 0,
  rating_support: 0,
  pros: '',
  cons: '',
  primary_use_case: '',
  reviewer_role: '',
  company_size: '',
  usage_context: ''
})

const errors = reactive({
  rating_overall: '',
  pros: '',
  cons: ''
})

function validate(): boolean {
  let valid = true
  errors.rating_overall = ''
  errors.pros = ''
  errors.cons = ''

  if (form.rating_overall < 1 || form.rating_overall > 5) {
    errors.rating_overall = 'الرجاء اختيار تقييم من 1 إلى 5'
    valid = false
  }

  if (!form.pros.trim()) {
    errors.pros = 'الرجاء تعبئة هذا الحقل'
    valid = false
  } else if (charCount(form.pros) > 500) {
    errors.pros = 'يجب ألا يتجاوز النص 500 حرف'
    valid = false
  }

  if (!form.cons.trim()) {
    errors.cons = 'الرجاء تعبئة هذا الحقل'
    valid = false
  } else if (charCount(form.cons) > 500) {
    errors.cons = 'يجب ألا يتجاوز النص 500 حرف'
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return

  submitting.value = true
  submitError.value = ''

  try {
    const payload: Record<string, any> = {
      rating_overall: form.rating_overall,
      pros: form.pros.trim(),
      cons: form.cons.trim()
    }

    // Only include optional fields if they have values
    if (form.rating_ease_of_use > 0) payload.rating_ease_of_use = form.rating_ease_of_use
    if (form.rating_value > 0) payload.rating_value = form.rating_value
    if (form.rating_accuracy > 0) payload.rating_accuracy = form.rating_accuracy
    if (form.rating_speed > 0) payload.rating_speed = form.rating_speed
    if (form.rating_support > 0) payload.rating_support = form.rating_support
    if (form.primary_use_case) payload.primary_use_case = form.primary_use_case
    if (form.reviewer_role) payload.reviewer_role = form.reviewer_role
    if (form.company_size) payload.company_size = form.company_size
    if (form.usage_context.trim()) payload.usage_context = form.usage_context.trim()

    await apiClient.post(`/tools/${props.toolSlug}/reviews`, payload)
    emit('success')
  } catch (err: any) {
    if (err.message?.includes('already reviewed')) {
      submitError.value = 'لقد قمت بمراجعة هذه الأداة مسبقاً'
    } else if (err.message?.includes('validation')) {
      submitError.value = 'الرجاء التأكد من صحة البيانات المدخلة'
    } else {
      submitError.value = 'حدث خطأ أثناء نشر المراجعة. الرجاء المحاولة مرة أخرى.'
    }
  } finally {
    submitting.value = false
  }
}

function handleClose() {
  emit('close')
}

// Handle Escape key to close modal
function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    handleClose()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown)
})
</script>

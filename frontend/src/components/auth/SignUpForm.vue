<template>
  <form @submit.prevent="handleSubmit" class="space-y-4" dir="rtl">
    <!-- Display Name -->
    <div>
      <label for="displayName" class="block text-sm font-medium text-gray-300 mb-1">
        الاسم
      </label>
      <input
        id="displayName"
        v-model="form.displayName"
        type="text"
        required
        :disabled="loading"
        class="w-full px-4 py-3 bg-dark-bg border border-gray-700 rounded-lg
               text-white placeholder-gray-500 focus:border-neon-blue focus:outline-none
               disabled:opacity-50"
        placeholder="اسمك"
      />
    </div>

    <!-- Email -->
    <div>
      <label for="email" class="block text-sm font-medium text-gray-300 mb-1">
        البريد الإلكتروني
      </label>
      <input
        id="email"
        v-model="form.email"
        type="email"
        required
        :disabled="loading"
        class="w-full px-4 py-3 bg-dark-bg border border-gray-700 rounded-lg
               text-white placeholder-gray-500 focus:border-neon-blue focus:outline-none
               disabled:opacity-50"
        placeholder="example@email.com"
      />
      <p v-if="errors.email" class="mt-1 text-red-400 text-sm">{{ errors.email }}</p>
    </div>

    <!-- Password -->
    <div>
      <label for="password" class="block text-sm font-medium text-gray-300 mb-1">
        كلمة المرور
      </label>
      <input
        id="password"
        v-model="form.password"
        type="password"
        required
        :disabled="loading"
        class="w-full px-4 py-3 bg-dark-bg border border-gray-700 rounded-lg
               text-white placeholder-gray-500 focus:border-neon-blue focus:outline-none
               disabled:opacity-50"
        placeholder="8 أحرف على الأقل"
      />
      <p v-if="errors.password" class="mt-1 text-red-400 text-sm">{{ errors.password }}</p>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="p-3 bg-red-500/10 border border-red-500/50 rounded-lg text-red-400 text-sm">
      {{ error }}
    </div>

    <!-- Submit Button -->
    <button
      type="submit"
      :disabled="loading"
      class="w-full py-3 bg-neon-blue text-dark-bg font-medium rounded-lg
             hover:bg-neon-blue-hover transition-colors disabled:opacity-50
             flex items-center justify-center gap-2"
    >
      <span v-if="loading" class="w-5 h-5 border-2 border-dark-bg border-t-transparent rounded-full animate-spin"></span>
      <span>{{ loading ? 'جاري إنشاء الحساب...' : 'إنشاء حساب' }}</span>
    </button>

    <!-- Switch to Sign In -->
    <p class="text-center text-gray-400 text-sm">
      لديك حساب؟
      <button
        type="button"
        @click="$emit('switch')"
        class="text-neon-blue hover:underline"
      >
        سجّل الدخول
      </button>
    </p>
  </form>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useSessionStore } from '@/stores/session'

const emit = defineEmits<{
  (e: 'switch'): void
  (e: 'success'): void
}>()

const sessionStore = useSessionStore()

const form = reactive({
  displayName: '',
  email: '',
  password: '',
})

const errors = reactive({
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref<string | null>(null)

function validate(): boolean {
  errors.email = ''
  errors.password = ''

  let valid = true

  // Email validation
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(form.email)) {
    errors.email = 'يرجى إدخال بريد إلكتروني صحيح'
    valid = false
  }

  // Password validation
  if (form.password.length < 8) {
    errors.password = 'كلمة المرور يجب أن تكون 8 أحرف على الأقل'
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return

  error.value = null
  loading.value = true

  try {
    await sessionStore.register(form.email, form.password, form.displayName)
    emit('success')
  } catch (err: any) {
    if (err?.response?.status === 409) {
      error.value = 'البريد الإلكتروني مستخدم بالفعل'
    } else if (err?.response?.status === 422) {
      error.value = 'يرجى التحقق من البيانات المدخلة'
    } else {
      error.value = 'حدث خطأ أثناء إنشاء الحساب. يرجى المحاولة مرة أخرى'
    }
  } finally {
    loading.value = false
  }
}
</script>

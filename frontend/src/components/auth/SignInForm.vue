<template>
  <form @submit.prevent="handleSubmit" class="space-y-4" dir="rtl">
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
        placeholder="••••••••"
      />
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
      <span>{{ loading ? 'جاري تسجيل الدخول...' : 'تسجيل الدخول' }}</span>
    </button>

    <!-- Switch to Sign Up -->
    <p class="text-center text-gray-400 text-sm">
      ليس لديك حساب؟
      <button
        type="button"
        @click="$emit('switch')"
        class="text-neon-blue hover:underline"
      >
        أنشئ حساباً
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
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref<string | null>(null)

async function handleSubmit() {
  error.value = null
  loading.value = true

  try {
    await sessionStore.login(form.email, form.password)
    emit('success')
  } catch (err: any) {
    if (err?.response?.status === 401) {
      error.value = 'البريد الإلكتروني أو كلمة المرور غير صحيحة'
    } else {
      error.value = 'حدث خطأ أثناء تسجيل الدخول. يرجى المحاولة مرة أخرى'
    }
  } finally {
    loading.value = false
  }
}
</script>

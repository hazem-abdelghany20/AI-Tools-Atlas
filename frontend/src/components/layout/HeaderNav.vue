<template>
  <header class="header-nav bg-dark-surface border-b border-gray-800 sticky top-0 z-40">
    <div class="container mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <!-- Logo (RTL: on the right) -->
        <div class="logo">
          <router-link to="/" class="text-xl font-bold text-neon-blue">
            AI Tools Atlas
          </router-link>
        </div>

        <!-- Navigation (Center) -->
        <nav class="hidden md:flex space-x-reverse space-x-6">
          <router-link to="/" class="nav-link">الرئيسية</router-link>
          <router-link to="/search" class="nav-link">تصفح</router-link>
          <router-link to="/bookmarks" class="nav-link">المفضلة</router-link>
        </nav>

        <!-- User menu (RTL: on the left) -->
        <div class="user-menu">
          <!-- Authenticated State -->
          <div v-if="isAuthenticated" class="relative">
            <button
              @click="showDropdown = !showDropdown"
              class="flex items-center gap-2 px-3 py-2 rounded-lg hover:bg-gray-800 transition-colors"
            >
              <!-- Avatar -->
              <div class="w-8 h-8 bg-neon-blue/20 rounded-full flex items-center justify-center">
                <span class="text-neon-blue font-medium text-sm">
                  {{ userInitials }}
                </span>
              </div>
              <span class="hidden sm:block text-white text-sm">{{ user?.display_name }}</span>
              <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </button>

            <!-- Dropdown Menu -->
            <Transition name="dropdown">
              <div
                v-if="showDropdown"
                class="absolute left-0 mt-2 w-48 bg-dark-surface border border-gray-800 rounded-lg shadow-xl z-50"
                dir="rtl"
              >
                <div class="p-2">
                  <router-link
                    to="/profile"
                    class="flex items-center gap-3 px-3 py-2 text-gray-300 hover:bg-gray-800 rounded-lg transition-colors"
                    @click="showDropdown = false"
                  >
                    <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                    <span>حسابي</span>
                  </router-link>
                  <router-link
                    to="/bookmarks"
                    class="flex items-center gap-3 px-3 py-2 text-gray-300 hover:bg-gray-800 rounded-lg transition-colors"
                    @click="showDropdown = false"
                  >
                    <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                    </svg>
                    <span>المفضلة</span>
                  </router-link>
                  <!-- Admin Links (only for admins) -->
                  <template v-if="user?.role === 'admin'">
                    <div class="px-3 py-1 text-xs text-gray-500 uppercase">Admin</div>
                    <router-link
                      to="/admin/tools"
                      class="flex items-center gap-3 px-3 py-2 text-neon-blue hover:bg-gray-800 rounded-lg transition-colors"
                      @click="showDropdown = false"
                    >
                      <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                      </svg>
                      <span>Tools</span>
                    </router-link>
                    <router-link
                      to="/admin/categories"
                      class="flex items-center gap-3 px-3 py-2 text-neon-blue hover:bg-gray-800 rounded-lg transition-colors"
                      @click="showDropdown = false"
                    >
                      <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                      </svg>
                      <span>Categories</span>
                    </router-link>
                    <router-link
                      to="/admin/tags"
                      class="flex items-center gap-3 px-3 py-2 text-neon-blue hover:bg-gray-800 rounded-lg transition-colors"
                      @click="showDropdown = false"
                    >
                      <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                      </svg>
                      <span>Tags</span>
                    </router-link>
                    <router-link
                      to="/admin/analytics"
                      class="flex items-center gap-3 px-3 py-2 text-neon-blue hover:bg-gray-800 rounded-lg transition-colors"
                      @click="showDropdown = false"
                    >
                      <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                      </svg>
                      <span>Analytics</span>
                    </router-link>
                    <router-link
                      to="/moderation/queue"
                      class="flex items-center gap-3 px-3 py-2 text-yellow-400 hover:bg-gray-800 rounded-lg transition-colors"
                      @click="showDropdown = false"
                    >
                      <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                      </svg>
                      <span>Moderation</span>
                    </router-link>
                  </template>
                  <div class="border-t border-gray-800 my-2"></div>
                  <button
                    @click="handleLogout"
                    class="w-full flex items-center gap-3 px-3 py-2 text-red-400 hover:bg-gray-800 rounded-lg transition-colors"
                  >
                    <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    <span>تسجيل الخروج</span>
                  </button>
                </div>
              </div>
            </Transition>
          </div>

          <!-- Guest State -->
          <div v-else class="flex items-center gap-3">
            <button
              @click="openAuthModal('signin')"
              class="px-4 py-2 text-gray-300 hover:text-white text-sm transition-colors"
            >
              تسجيل الدخول
            </button>
            <button
              @click="openAuthModal('signup')"
              class="px-4 py-2 bg-neon-blue text-dark-bg font-medium text-sm rounded-lg
                     hover:bg-neon-blue-hover transition-colors"
            >
              إنشاء حساب
            </button>
          </div>
        </div>

        <!-- Mobile menu button -->
        <button
          class="md:hidden p-2 text-gray-400 hover:text-white"
          @click="toggleMobileMenu"
        >
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Mobile menu -->
      <Transition name="mobile-menu">
        <div v-if="mobileMenuOpen" class="md:hidden py-4 border-t border-gray-800">
          <div class="flex flex-col space-y-2" dir="rtl">
            <router-link
              to="/"
              class="px-4 py-2 text-gray-300 hover:text-white hover:bg-gray-800 rounded-lg"
              @click="mobileMenuOpen = false"
            >
              الرئيسية
            </router-link>
            <router-link
              to="/search"
              class="px-4 py-2 text-gray-300 hover:text-white hover:bg-gray-800 rounded-lg"
              @click="mobileMenuOpen = false"
            >
              تصفح
            </router-link>
            <router-link
              to="/bookmarks"
              class="px-4 py-2 text-gray-300 hover:text-white hover:bg-gray-800 rounded-lg"
              @click="mobileMenuOpen = false"
            >
              المفضلة
            </router-link>

            <!-- Mobile Auth -->
            <div v-if="!isAuthenticated" class="pt-4 border-t border-gray-800 space-y-2">
              <button
                @click="openAuthModal('signin'); mobileMenuOpen = false"
                class="w-full px-4 py-2 text-gray-300 hover:text-white text-center hover:bg-gray-800 rounded-lg"
              >
                تسجيل الدخول
              </button>
              <button
                @click="openAuthModal('signup'); mobileMenuOpen = false"
                class="w-full px-4 py-2 bg-neon-blue text-dark-bg font-medium rounded-lg"
              >
                إنشاء حساب
              </button>
            </div>
            <div v-else class="pt-4 border-t border-gray-800">
              <button
                @click="handleLogout; mobileMenuOpen = false"
                class="w-full px-4 py-2 text-red-400 hover:bg-gray-800 rounded-lg text-right"
              >
                تسجيل الخروج
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </div>
  </header>

  <!-- Auth Modal -->
  <AuthModal
    :is-open="authModalOpen"
    :initial-mode="authModalMode"
    @close="authModalOpen = false"
  />
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useSessionStore } from '@/stores/session'
import AuthModal from '@/components/auth/AuthModal.vue'

const sessionStore = useSessionStore()

const mobileMenuOpen = ref(false)
const showDropdown = ref(false)
const authModalOpen = ref(false)
const authModalMode = ref<'signin' | 'signup'>('signin')

const isAuthenticated = computed(() => sessionStore.isAuthenticated)
const user = computed(() => sessionStore.user)

const userInitials = computed(() => {
  if (!user.value?.display_name) return '?'
  return user.value.display_name.charAt(0).toUpperCase()
})

function toggleMobileMenu() {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

function openAuthModal(mode: 'signin' | 'signup') {
  authModalMode.value = mode
  authModalOpen.value = true
}

function handleLogout() {
  showDropdown.value = false
  sessionStore.logout()
}

// Close dropdown when clicking outside
function handleClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (!target.closest('.user-menu')) {
    showDropdown.value = false
  }
}

// Fetch current user on mount
onMounted(() => {
  sessionStore.fetchCurrentUser()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.nav-link {
  @apply text-gray-300 hover:text-white transition-colors;
}

.router-link-active.nav-link {
  @apply text-neon-blue;
}

/* Dropdown Transition */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Mobile Menu Transition */
.mobile-menu-enter-active,
.mobile-menu-leave-active {
  transition: all 0.2s ease;
}

.mobile-menu-enter-from,
.mobile-menu-leave-to {
  opacity: 0;
  max-height: 0;
}
</style>

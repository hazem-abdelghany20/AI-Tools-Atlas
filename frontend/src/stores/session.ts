import { defineStore } from 'pinia'
import { ref } from 'vue'
import { apiClient } from '../lib/apiClient'
import type { User } from '../lib/types'

export const useSessionStore = defineStore('session', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = ref(false)
  const isLoading = ref(false)

  async function fetchCurrentUser() {
    isLoading.value = true
    try {
      const response = await apiClient.get<{ data: User }>('/me')
      user.value = response.data
      isAuthenticated.value = true
    } catch (error) {
      user.value = null
      isAuthenticated.value = false
    } finally {
      isLoading.value = false
    }
  }

  async function login(email: string, password: string): Promise<User> {
    isLoading.value = true
    try {
      const response = await apiClient.post<{ data: { user: User } }>('/auth/login', {
        email,
        password,
      })
      user.value = response.data.user
      isAuthenticated.value = true
      return response.data.user
    } catch (error) {
      throw error
    } finally {
      isLoading.value = false
    }
  }

  async function register(email: string, password: string, displayName: string): Promise<User> {
    isLoading.value = true
    try {
      const response = await apiClient.post<{ data: { user: User } }>('/auth/register', {
        email,
        password,
        display_name: displayName,
      })
      user.value = response.data.user
      isAuthenticated.value = true
      return response.data.user
    } catch (error) {
      throw error
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    isLoading.value = true
    try {
      await apiClient.post('/auth/logout')
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      user.value = null
      isAuthenticated.value = false
      isLoading.value = false
    }
  }

  return {
    user,
    isAuthenticated,
    isLoading,
    fetchCurrentUser,
    login,
    register,
    logout,
  }
})

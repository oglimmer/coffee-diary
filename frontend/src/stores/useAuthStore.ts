import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { authService } from '@/services/authService'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isAuthenticated = computed(() => user.value !== null)

  async function login(username: string, password: string) {
    user.value = await authService.login(username, password)
  }

  async function register(username: string, password: string) {
    user.value = await authService.register(username, password)
  }

  async function logout() {
    await authService.logout()
    user.value = null
  }

  async function fetchUser() {
    try {
      user.value = await authService.getMe()
    } catch {
      user.value = null
    }
  }

  return { user, isAuthenticated, login, register, logout, fetchUser }
})

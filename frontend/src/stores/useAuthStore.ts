import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { authService } from '@/services/authService'
import type { User } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const initialized = ref(false)
  const restoring = ref(false)
  const isAuthenticated = computed(() => user.value !== null)
  const isReady = computed(() => initialized.value)

  let restorePromise: Promise<void> | null = null

  async function login(username: string, password: string) {
    user.value = await authService.login(username, password)
    initialized.value = true
  }

  async function register(username: string, password: string) {
    user.value = await authService.register(username, password)
    initialized.value = true
  }

  async function logout() {
    await authService.logout()
    user.value = null
    initialized.value = true
  }

  async function fetchUser() {
    try {
      user.value = await authService.getMe()
    } catch {
      user.value = null
    } finally {
      initialized.value = true
    }
  }

  async function restoreSession() {
    if (initialized.value) {
      return
    }

    if (!restorePromise) {
      restoring.value = true
      restorePromise = fetchUser().finally(() => {
        restoring.value = false
        restorePromise = null
      })
    }

    await restorePromise
  }

  return { user, initialized, restoring, isAuthenticated, isReady, login, register, logout, fetchUser, restoreSession }
})

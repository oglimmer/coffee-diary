import api from './api'
import type { User } from '@/types'

export const authService = {
  logout(): void {
    window.location.href = '/api/auth/logout'
  },

  async getMe(): Promise<User> {
    const { data } = await api.get<User>('/auth/me')
    return data
  },

  redirectToLogin(): void {
    window.location.href = '/api/auth/login'
  },
}

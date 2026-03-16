import api from './api'
import type { User } from '@/types'

export const authService = {
  async login(username: string, password: string): Promise<User> {
    const { data } = await api.post<User>('/auth/login', { username, password })
    return data
  },

  async register(username: string, password: string): Promise<User> {
    const { data } = await api.post<User>('/auth/register', { username, password })
    return data
  },

  async logout(): Promise<void> {
    await api.post('/auth/logout')
  },

  async getMe(): Promise<User> {
    const { data } = await api.get<User>('/auth/me')
    return data
  },
}

import api from './api'
import type { Sieve } from '@/types'

export const sieveService = {
  async getSieves(): Promise<Sieve[]> {
    const { data } = await api.get<Sieve[]>('/sieves')
    return data
  },

  async createSieve(name: string): Promise<Sieve> {
    const { data } = await api.post<Sieve>('/sieves', { name })
    return data
  },

  async deleteSieve(id: number): Promise<void> {
    await api.delete(`/sieves/${id}`)
  },
}

import api from './api'
import type { Coffee } from '@/types'

export const coffeeService = {
  async getCoffees(): Promise<Coffee[]> {
    const { data } = await api.get<Coffee[]>('/coffees')
    return data
  },

  async createCoffee(name: string): Promise<Coffee> {
    const { data } = await api.post<Coffee>('/coffees', { name })
    return data
  },

  async deleteCoffee(id: number): Promise<void> {
    await api.delete(`/coffees/${id}`)
  },
}

import api from './api'
import type { DiaryEntry, DiaryEntryRequest, Page } from '@/types'

export interface DiaryQueryParams {
  page?: number
  size?: number
  sort?: string
  dateFrom?: string
  dateTo?: string
  coffeeId?: number
  sieveId?: number
  minRating?: number
}

export const diaryService = {
  async getEntries(params: DiaryQueryParams = {}): Promise<Page<DiaryEntry>> {
    const { data } = await api.get<Page<DiaryEntry>>('/diary-entries', { params })
    return data
  },

  async getEntry(id: number): Promise<DiaryEntry> {
    const { data } = await api.get<DiaryEntry>(`/diary-entries/${id}`)
    return data
  },

  async createEntry(entry: DiaryEntryRequest): Promise<DiaryEntry> {
    const { data } = await api.post<DiaryEntry>('/diary-entries', entry)
    return data
  },

  async updateEntry(id: number, entry: DiaryEntryRequest): Promise<DiaryEntry> {
    const { data } = await api.put<DiaryEntry>(`/diary-entries/${id}`, entry)
    return data
  },

  async deleteEntry(id: number): Promise<void> {
    await api.delete(`/diary-entries/${id}`)
  },
}

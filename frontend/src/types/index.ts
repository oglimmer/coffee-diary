export interface User {
  id: number
  username: string
}

export interface Sieve {
  id: number
  name: string
}

export interface Coffee {
  id: number
  name: string
}

export interface DiaryEntry {
  id: number
  dateTime: string
  sieveId: number | null
  sieveName: string | null
  temperature: number
  coffeeId: number | null
  coffeeName: string | null
  grindSize: number
  inputWeight: number
  outputWeight: number
  timeSeconds: number
  rating: number
  notes: string
}

export interface DiaryEntryRequest {
  dateTime: string
  sieveId: number | null
  temperature: number
  coffeeId: number | null
  grindSize: number | null
  inputWeight: number | null
  outputWeight: number | null
  timeSeconds: number | null
  rating: number
  notes: string
}

export interface Page<T> {
  content: T[]
  totalElements: number
  totalPages: number
  number: number
  size: number
}

import axios from 'axios'
import router from '@/router'
import { useAuthStore } from '@/stores/useAuthStore'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const url = error.config?.url ?? ''
    if (error.response?.status === 401 && !url.includes('/auth/me')) {
      const auth = useAuthStore()
      auth.user = null
      router.push({ path: '/login', query: { expired: '1' } })
    }
    return Promise.reject(error)
  },
)

export default api

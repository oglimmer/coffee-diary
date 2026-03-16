import axios from 'axios'
import router from '@/router'

const api = axios.create({
  baseURL: '/api',
  withCredentials: true,
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    const url = error.config?.url ?? ''
    if (error.response?.status === 401 && !url.includes('/auth/me')) {
      router.push('/login')
    }
    return Promise.reject(error)
  },
)

export default api

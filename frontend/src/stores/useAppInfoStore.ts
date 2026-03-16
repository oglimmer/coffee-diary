import { ref } from 'vue'
import { defineStore } from 'pinia'
import { appInfoService, type AppInfo } from '@/services/appInfoService'

export const useAppInfoStore = defineStore('appInfo', () => {
  const info = ref<AppInfo | null>(null)

  async function load() {
    try {
      info.value = await appInfoService.fetchAppInfo()
    } catch {
      info.value = null
    }
  }

  return { info, load }
})

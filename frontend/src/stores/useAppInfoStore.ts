import { ref } from 'vue'
import { defineStore } from 'pinia'
import { appInfoService, type AppInfo } from '@/services/appInfoService'

declare const __APP_VERSION__: string
declare const __BUILD_TIME__: string
declare const __GIT_COMMIT__: string

export interface CombinedAppInfo {
  frontend: { version: string; buildTime: string; gitCommit: string }
  backend: AppInfo | null
}

export const useAppInfoStore = defineStore('appInfo', () => {
  const info = ref<AppInfo | null>(null)
  const combined = ref<CombinedAppInfo>({
    frontend: { version: __APP_VERSION__, buildTime: __BUILD_TIME__, gitCommit: __GIT_COMMIT__ },
    backend: null,
  })

  async function load() {
    try {
      info.value = await appInfoService.fetchAppInfo()
      combined.value.backend = info.value
    } catch {
      info.value = null
      combined.value.backend = null
    }
  }

  return { info, combined, load }
})

export interface AppInfo {
  app: {
    name: string
    version: string
  }
  build: {
    time: string
  }
  git: {
    commit: string
  }
}

export const appInfoService = {
  async fetchAppInfo(): Promise<AppInfo> {
    const response = await fetch('/actuator/info')
    return response.json()
  },
}

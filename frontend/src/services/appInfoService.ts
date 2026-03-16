import axios from 'axios'

export interface AppInfo {
  app: {
    name: string
    version: string
  }
  build: {
    time: string
  }
}

export const appInfoService = {
  async fetchAppInfo(): Promise<AppInfo> {
    const { data } = await axios.get<AppInfo>('/actuator/info')
    return data
  },
}

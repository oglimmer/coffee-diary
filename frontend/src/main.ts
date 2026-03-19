import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import { useAuthStore } from '@/stores/useAuthStore'
import '@fontsource-variable/fraunces'
import '@fontsource-variable/outfit'
import './assets/main.css'

const app = createApp(App)
app.use(createPinia())

// Restore session before router processes any route
const auth = useAuthStore()
auth.fetchUser().finally(() => {
  app.use(router)
  app.mount('#app')
})

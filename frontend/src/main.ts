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
app.use(router)

const auth = useAuthStore()
auth.restoreSession()

app.mount('#app')

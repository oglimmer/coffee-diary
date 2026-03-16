<script setup lang="ts">
import { useAuthStore } from '@/stores/useAuthStore'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

async function handleLogout() {
  await auth.logout()
  router.push('/login')
}
</script>

<template>
  <header class="app-header">
    <div class="header-inner">
      <RouterLink to="/" class="logo">
        <span class="logo-text">Coffee<br/>Diary</span>
      </RouterLink>
      <div v-if="auth.isAuthenticated" class="header-right">
        <span class="username">{{ auth.user?.username }}</span>
        <button class="btn btn-outline btn-sm" @click="handleLogout">Logout</button>
      </div>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  border-bottom: 1px solid var(--border);
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--bg);
  backdrop-filter: blur(12px);
}

.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  text-decoration: none;
  color: var(--text);
}

.logo:hover {
  color: var(--text);
}

.logo-text {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 400;
  line-height: 1.1;
  letter-spacing: -0.02em;
  font-variation-settings: 'SOFT' 100, 'WONK' 1;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.username {
  color: var(--text-muted);
  font-size: 13px;
  font-weight: 400;
  letter-spacing: 0.02em;
}
</style>

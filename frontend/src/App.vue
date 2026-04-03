<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { useAppInfoStore } from '@/stores/useAppInfoStore'
const auth = useAuthStore()
const appInfo = useAppInfoStore()
const showBootSplash = computed(() => auth.restoring && !auth.initialized)

onMounted(() => {
  appInfo.load()
})
</script>

<template>
  <div v-if="showBootSplash" class="boot-splash">
    <div class="loading-dots" aria-label="Loading">
      <div class="loading-dot"></div>
      <div class="loading-dot"></div>
      <div class="loading-dot"></div>
    </div>
  </div>
  <RouterView v-else />
  <footer class="app-footer">
    <nav class="footer-nav">
      <RouterLink to="/privacy">Privacy Policy</RouterLink>
      <span class="footer-sep">&middot;</span>
      <RouterLink to="/terms">Terms &amp; Conditions</RouterLink>
      <span class="footer-sep">&middot;</span>
      <RouterLink to="/imprint">Imprint</RouterLink>
      <span class="footer-sep">&middot;</span>
      <RouterLink to="/developer">Developer</RouterLink>
    </nav>
    <div v-if="appInfo.info?.app" class="footer-build">
      {{ appInfo.info.app.name }} v{{ appInfo.info.app.version }}
      <span class="footer-sep">&middot;</span>
      Built {{ appInfo.info.build.time }}
    </div>
  </footer>
</template>

<style scoped>
.boot-splash {
  min-height: 100vh;
  display: grid;
  place-items: center;
}

.loading-dots {
  display: flex;
  gap: 8px;
}

.loading-dot {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: var(--accent);
  animation: pulse 1s ease-in-out infinite;
}

.loading-dot:nth-child(2) {
  animation-delay: 0.15s;
}

.loading-dot:nth-child(3) {
  animation-delay: 0.3s;
}

.app-footer {
  text-align: center;
  padding: 24px 24px 20px;
  border-top: 1px solid var(--border);
  margin-top: auto;
}

.footer-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-size: 12px;
  letter-spacing: 0.02em;
  margin-bottom: 12px;
}

.footer-nav a {
  color: var(--text-dim);
  text-decoration: none;
  transition: color 0.2s ease;
}

.footer-nav a:hover {
  color: var(--accent);
}

.footer-nav a.router-link-active {
  color: var(--text-muted);
  pointer-events: none;
}

.footer-build {
  font-size: 11px;
  color: var(--text-dim);
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.footer-sep {
  margin: 0 4px;
  opacity: 0.4;
  color: var(--text-dim);
}

@keyframes pulse {
  0%, 80%, 100% {
    opacity: 0.25;
    transform: translateY(0);
  }
  40% {
    opacity: 1;
    transform: translateY(-4px);
  }
}
</style>

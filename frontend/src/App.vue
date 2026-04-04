<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { useAppInfoStore } from '@/stores/useAppInfoStore'
const auth = useAuthStore()
const appInfo = useAppInfoStore()
const showBootSplash = computed(() => auth.restoring && !auth.initialized)
const showBuildInfo = ref(false)

function closeBuildInfo() {
  showBuildInfo.value = false
}

onMounted(() => {
  appInfo.load()
  document.addEventListener('click', closeBuildInfo)
})

onUnmounted(() => {
  document.removeEventListener('click', closeBuildInfo)
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
      <span class="build-info-wrapper">
        <button class="build-info-btn" @click.stop="showBuildInfo = !showBuildInfo" aria-label="Build info">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/></svg>
        </button>
        <div v-if="showBuildInfo" class="build-info-popup" @click.stop>
          <div class="build-info-row">
            <span class="build-info-label">Frontend</span>
            <span>v{{ appInfo.combined.frontend.version }}</span>
            <span class="build-info-dim">#{{ appInfo.combined.frontend.gitCommit }}</span>
            <span class="build-info-dim">{{ appInfo.combined.frontend.buildTime }}</span>
          </div>
          <div v-if="appInfo.combined.backend" class="build-info-row">
            <span class="build-info-label">Backend</span>
            <span>v{{ appInfo.combined.backend.app.version }}</span>
            <span class="build-info-dim">#{{ appInfo.combined.backend.git.commit }}</span>
            <span class="build-info-dim">{{ appInfo.combined.backend.build.time }}</span>
          </div>
        </div>
      </span>
    </nav>
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

.build-info-wrapper {
  position: relative;
  display: inline-flex;
  align-items: center;
}

.build-info-btn {
  background: none;
  border: none;
  color: var(--text-dim);
  cursor: pointer;
  padding: 2px;
  display: inline-flex;
  align-items: center;
  transition: color 0.2s ease;
  line-height: 1;
}

.build-info-btn:hover {
  color: var(--accent);
}

.build-info-popup {
  position: absolute;
  bottom: calc(100% + 8px);
  right: -8px;
  background: var(--surface, #1a1a1a);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 10px 14px;
  font-size: 11px;
  white-space: nowrap;
  z-index: 100;
  display: flex;
  flex-direction: column;
  gap: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.build-info-row {
  display: flex;
  gap: 8px;
  align-items: center;
  color: var(--text-muted);
}

.build-info-label {
  font-weight: 600;
  color: var(--text);
  min-width: 60px;
}

.build-info-dim {
  color: var(--text-dim);
  font-size: 10px;
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

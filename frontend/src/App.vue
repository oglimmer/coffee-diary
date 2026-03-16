<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useAppInfoStore } from '@/stores/useAppInfoStore'

const appInfo = useAppInfoStore()

onMounted(() => {
  appInfo.load()
})
</script>

<template>
  <RouterView />
  <footer v-if="appInfo.info?.app" class="app-footer">
    <span>{{ appInfo.info.app.name }} v{{ appInfo.info.app.version }}</span>
    <span class="footer-sep">&middot;</span>
    <span>Built {{ appInfo.info.build.time }}</span>
  </footer>
</template>

<style scoped>
.app-footer {
  text-align: center;
  padding: 20px 24px;
  font-size: 11px;
  color: var(--text-dim);
  letter-spacing: 0.04em;
  text-transform: uppercase;
  border-top: 1px solid var(--border);
  margin-top: auto;
}

.footer-sep {
  margin: 0 8px;
  opacity: 0.4;
}
</style>

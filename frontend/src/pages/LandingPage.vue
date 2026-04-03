<script setup lang="ts">
import { useAuthStore } from '@/stores/useAuthStore'

const auth = useAuthStore()

function handleLogin() {
  auth.redirectToLogin()
}
</script>

<template>
  <div class="landing-page">
    <div class="landing-ambient"></div>
    <div class="landing-card fade-up">
      <h1 class="landing-title">Coffee<br/>Diary</h1>
      <p class="landing-subtitle">Track your espresso journey.<br/>Dial in the perfect shot.</p>

      <button class="btn btn-primary btn-block" @click="handleLogin">
        Sign in with SSO
      </button>

      <p class="landing-hint">You will be redirected to our identity provider</p>
    </div>

  </div>
</template>

<style scoped>
.landing-page {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

.landing-ambient {
  position: absolute;
  width: 600px;
  height: 600px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--accent-glow-strong) 0%, transparent 70%);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  pointer-events: none;
  filter: blur(80px);
}

.landing-card {
  background: var(--surface);
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-xl);
  padding: 56px 44px;
  width: 100%;
  max-width: 400px;
  text-align: center;
  position: relative;
  z-index: 1;
}

.landing-title {
  font-size: 36px;
  font-family: var(--font-display);
  font-weight: 400;
  color: var(--text);
  margin: 0 0 8px;
  line-height: 1;
  letter-spacing: -0.02em;
  font-variation-settings: 'SOFT' 100, 'WONK' 1;
}

.landing-subtitle {
  color: var(--text-muted);
  margin: 0 0 40px;
  font-size: 14px;
  letter-spacing: 0.02em;
  line-height: 1.6;
}

.landing-hint {
  margin-top: 16px;
  font-size: 12px;
  color: var(--text-dim);
  letter-spacing: 0.02em;
}
</style>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import LegalNav from '@/components/LegalNav.vue'

const auth = useAuthStore()
const router = useRouter()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''

  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match.'
    return
  }

  if (password.value.length < 4) {
    error.value = 'Password must be at least 4 characters.'
    return
  }

  loading.value = true
  try {
    await auth.register(username.value, password.value)
    router.push('/')
  } catch (e: unknown) {
    const err = e as { response?: { data?: { message?: string } } }
    error.value = err.response?.data?.message || 'Registration failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-ambient"></div>
    <div class="auth-card fade-up">
      <h1 class="auth-title">Coffee<br/>Diary</h1>
      <p class="auth-subtitle">Create a new account</p>

      <form @submit.prevent="handleSubmit" class="auth-form">
        <div v-if="error" class="alert alert-error">{{ error }}</div>

        <div class="form-group">
          <label for="username">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            autocomplete="username"
            placeholder="Choose a username"
          />
        </div>

        <div class="form-group">
          <label for="password">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            autocomplete="new-password"
            placeholder="Choose a password"
          />
        </div>

        <div class="form-group">
          <label for="confirmPassword">Confirm Password</label>
          <input
            id="confirmPassword"
            v-model="confirmPassword"
            type="password"
            required
            autocomplete="new-password"
            placeholder="Repeat your password"
          />
        </div>

        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? 'Creating account...' : 'Create Account' }}
        </button>
      </form>

      <p class="auth-link">
        Already have an account? <RouterLink to="/login">Sign in</RouterLink>
      </p>
    </div>

    <LegalNav />
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

.auth-page :deep(.legal-nav) {
  position: absolute;
  bottom: 24px;
  border-top: none;
  margin-top: 0;
}

.auth-ambient {
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

.auth-card {
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

.auth-title {
  font-size: 36px;
  font-family: var(--font-display);
  font-weight: 400;
  color: var(--text);
  margin: 0 0 8px;
  line-height: 1;
  letter-spacing: -0.02em;
  font-variation-settings: 'SOFT' 100, 'WONK' 1;
}

.auth-subtitle {
  color: var(--text-dim);
  margin: 0 0 40px;
  font-size: 14px;
  letter-spacing: 0.02em;
}

.auth-form {
  text-align: left;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.auth-form .btn-block {
  margin-top: 8px;
}

.auth-link {
  margin-top: 28px;
  font-size: 13px;
  color: var(--text-dim);
}

.auth-link a {
  color: var(--accent);
  font-weight: 500;
}
</style>

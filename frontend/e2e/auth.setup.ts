import { test as setup, expect } from '@playwright/test'

const KEYCLOAK_USER = process.env.E2E_USER ?? 'test'
const KEYCLOAK_PASS = process.env.E2E_PASS ?? 'test'

setup('authenticate via Keycloak', async ({ page }) => {
  // Go to landing page, which triggers the OIDC flow
  await page.goto('/landing')
  await page.getByRole('button', { name: 'Sign in with SSO' }).click()

  // Now on the Keycloak login page
  await page.waitForURL(/id\.oglimmer\.de/)
  await page.locator('#username').fill(KEYCLOAK_USER)
  await page.locator('#password').fill(KEYCLOAK_PASS)
  await page.locator('#kc-login').click()

  // Redirected back to the app — wait for the diary page
  await page.waitForURL('http://localhost:5173/')
  await expect(page.locator('.username')).toBeVisible()

  // Persist session (cookies) for all tests
  await page.context().storageState({ path: 'e2e/.auth/state.json' })
})

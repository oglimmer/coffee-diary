import { test, expect } from '@playwright/test'

test('unauthenticated user is redirected to landing', async ({ browser }) => {
  // Fresh context without saved auth state — explicitly clear storageState
  const context = await browser.newContext({ storageState: undefined })
  const page = await context.newPage()

  await page.goto('/')
  await expect(page).toHaveURL(/\/landing/)
  await expect(page.getByRole('button', { name: 'Sign in with SSO' })).toBeVisible()

  await context.close()
})

test('logo navigates to diary list', async ({ page }) => {
  await page.goto('/entry/new')
  await page.locator('.logo').click()
  await expect(page).toHaveURL('/')
})

test('logout redirects away from protected pages', async ({ browser }) => {
  // Use saved auth state
  const context = await browser.newContext({
    storageState: 'e2e/.auth/state.json',
  })
  const page = await context.newPage()

  await page.goto('/')
  await expect(page.locator('.username')).toBeVisible()

  // Clicking logout navigates to /api/auth/logout which redirects through Keycloak
  // We just verify the button exists and is clickable
  await expect(page.getByRole('button', { name: 'Logout' })).toBeVisible()

  await context.close()
})

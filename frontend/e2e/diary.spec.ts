import { test, expect } from '@playwright/test'

test.describe('Diary list page', () => {
  test('shows the diary page with header', async ({ page }) => {
    await page.goto('/')
    await expect(page.locator('.page-title')).toHaveText('My Entries')
    await expect(page.locator('.username')).toBeVisible()
  })

  test('shows empty state or entries table', async ({ page }) => {
    await page.goto('/')
    const empty = page.locator('.empty-state')
    const table = page.locator('.entries-table')
    await expect(empty.or(table)).toBeVisible()
  })

  test('can navigate to new entry form', async ({ page }) => {
    await page.goto('/')
    await page.getByRole('link', { name: '+ New Entry' }).click()
    await expect(page).toHaveURL('/entry/new')
    await expect(page.locator('.form-title')).toHaveText('New Entry')
  })

  test('can toggle filters', async ({ page }) => {
    await page.goto('/')
    await page.getByRole('button', { name: 'Filters' }).click()
    await expect(page.locator('.filter-bar')).toBeVisible()
    await page.getByRole('button', { name: 'Hide Filters' }).click()
    await expect(page.locator('.filter-bar')).not.toBeVisible()
  })
})

test.describe('Diary entry form', () => {
  test('can create and delete an entry', async ({ page }) => {
    await page.goto('/entry/new')

    // Fill minimal form fields
    await page.locator('#temperature').fill('93')
    await page.locator('#grindSize').fill('2.5')
    await page.locator('#inputWeight').fill('18')
    await page.locator('#outputWeight').fill('36')
    await page.locator('#timeSeconds').fill('25')

    await page.getByRole('button', { name: 'Save Entry' }).click()

    // Should redirect to diary list
    await expect(page).toHaveURL('/')

    // The entry should appear in the table
    const firstRow = page.locator('.entry-row').first()
    await expect(firstRow).toBeVisible()

    // Clean up: delete the entry we just created
    await firstRow.locator('.btn-icon').click()
    await page.locator('.dialog-box').getByRole('button', { name: 'Delete' }).click()
  })

  test('cancel returns to diary list', async ({ page }) => {
    await page.goto('/entry/new')
    await page.getByRole('button', { name: 'Cancel' }).click()
    await expect(page).toHaveURL('/')
  })
})

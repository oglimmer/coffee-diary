<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { diaryService, type DiaryQueryParams } from '@/services/diaryService'
import { coffeeService } from '@/services/coffeeService'
import { sieveService } from '@/services/sieveService'
import type { DiaryEntry, Coffee, Sieve, Page } from '@/types'
import AppHeader from '@/components/AppHeader.vue'
import StarRating from '@/components/StarRating.vue'
import PaginationControls from '@/components/PaginationControls.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'

const router = useRouter()

const entries = ref<DiaryEntry[]>([])
const page = ref(0)
const totalPages = ref(0)
const totalElements = ref(0)
const loading = ref(false)

// Filters
const dateFrom = ref('')
const dateTo = ref('')
const coffeeFilter = ref<number | ''>('')
const sieveFilter = ref<number | ''>('')
const minRating = ref<number | ''>('')
const showFilters = ref(false)

// Dropdown data
const coffees = ref<Coffee[]>([])
const sieves = ref<Sieve[]>([])

// Delete dialog
const deleteDialogVisible = ref(false)
const entryToDelete = ref<DiaryEntry | null>(null)

async function loadEntries() {
  loading.value = true
  try {
    const params: DiaryQueryParams = {
      page: page.value,
      size: 15,
      sort: 'dateTime,desc',
    }
    if (dateFrom.value) params.dateFrom = dateFrom.value
    if (dateTo.value) params.dateTo = dateTo.value
    if (coffeeFilter.value !== '') params.coffeeId = coffeeFilter.value as number
    if (sieveFilter.value !== '') params.sieveId = sieveFilter.value as number
    if (minRating.value !== '') params.minRating = minRating.value as number

    const result: Page<DiaryEntry> = await diaryService.getEntries(params)
    entries.value = result.content
    totalPages.value = result.totalPages
    totalElements.value = result.totalElements
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
}

async function loadDropdowns() {
  const [c, s] = await Promise.all([coffeeService.getCoffees(), sieveService.getSieves()])
  coffees.value = c
  sieves.value = s
}

function applyFilters() {
  page.value = 0
  loadEntries()
}

function clearFilters() {
  dateFrom.value = ''
  dateTo.value = ''
  coffeeFilter.value = ''
  sieveFilter.value = ''
  minRating.value = ''
  page.value = 0
  loadEntries()
}

function editEntry(entry: DiaryEntry) {
  router.push(`/entry/${entry.id}/edit`)
}

function confirmDelete(entry: DiaryEntry) {
  entryToDelete.value = entry
  deleteDialogVisible.value = true
}

async function handleDelete() {
  if (entryToDelete.value) {
    await diaryService.deleteEntry(entryToDelete.value.id)
    deleteDialogVisible.value = false
    entryToDelete.value = null
    loadEntries()
  }
}

function formatDate(dt: string) {
  const d = new Date(dt)
  return d.toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
}

function formatTime(dt: string) {
  const d = new Date(dt)
  return d.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' })
}

function formatRatio(input: number, output: number) {
  if (!input || !output) return '-'
  return `${input}g \u2192 ${output}g`
}

function formatSeconds(s: number) {
  if (!s) return '-'
  return `${s}s`
}

watch(page, () => {
  loadEntries()
})

onMounted(() => {
  loadEntries()
  loadDropdowns()
})
</script>

<template>
  <div class="page-wrapper">
    <AppHeader />

    <main class="container fade-up">
      <div class="page-top">
        <div>
          <h2 class="page-title">My Entries</h2>
          <p class="page-count">{{ totalElements }} total entries</p>
        </div>
        <div class="page-top-actions">
          <button class="btn btn-outline btn-sm" @click="showFilters = !showFilters">
            {{ showFilters ? 'Hide Filters' : 'Filters' }}
          </button>
          <RouterLink to="/entry/new" class="btn btn-primary">+ New Entry</RouterLink>
        </div>
      </div>

      <!-- Filter Bar -->
      <Transition name="slide">
        <div v-if="showFilters" class="filter-bar">
          <div class="filter-row">
            <div class="filter-group">
              <label>From</label>
              <input type="date" v-model="dateFrom" />
            </div>
            <div class="filter-group">
              <label>To</label>
              <input type="date" v-model="dateTo" />
            </div>
            <div class="filter-group">
              <label>Coffee</label>
              <select v-model="coffeeFilter">
                <option value="">All</option>
                <option v-for="c in coffees" :key="c.id" :value="c.id">{{ c.name }}</option>
              </select>
            </div>
            <div class="filter-group">
              <label>Sieve</label>
              <select v-model="sieveFilter">
                <option value="">All</option>
                <option v-for="s in sieves" :key="s.id" :value="s.id">{{ s.name }}</option>
              </select>
            </div>
            <div class="filter-group">
              <label>Min Rating</label>
              <select v-model="minRating">
                <option value="">Any</option>
                <option v-for="n in 5" :key="n" :value="n">{{ n }}+</option>
              </select>
            </div>
            <div class="filter-actions">
              <button class="btn btn-sm btn-primary" @click="applyFilters">Apply</button>
              <button class="btn btn-sm btn-outline" @click="clearFilters">Clear</button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Entries Table -->
      <div class="entries-table-wrapper" v-if="!loading && entries.length > 0">
        <table class="entries-table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Coffee</th>
              <th>Sieve</th>
              <th>Temp</th>
              <th>Grind</th>
              <th>Ratio</th>
              <th>Time</th>
              <th>Rating</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="entry in entries" :key="entry.id" class="entry-row" @click="editEntry(entry)">
              <td>
                <div class="date-cell">
                  <span class="date-main">{{ formatDate(entry.dateTime) }}</span>
                  <span class="date-time">{{ formatTime(entry.dateTime) }}</span>
                </div>
              </td>
              <td class="cell-coffee">{{ entry.coffeeName || '-' }}</td>
              <td>{{ entry.sieveName || '-' }}</td>
              <td class="cell-mono">{{ entry.temperature ? entry.temperature + '\u00B0C' : '-' }}</td>
              <td class="cell-mono">{{ entry.grindSize || '-' }}</td>
              <td class="cell-mono">{{ formatRatio(entry.inputWeight, entry.outputWeight) }}</td>
              <td class="cell-mono">{{ formatSeconds(entry.timeSeconds) }}</td>
              <td>
                <StarRating :model-value="entry.rating" :readonly="true" :size="14" />
              </td>
              <td class="actions-cell" @click.stop>
                <button class="btn-icon" title="Delete" @click="confirmDelete(entry)">
                  <svg width="16" height="16" viewBox="0 0 16 16" fill="none"><path d="M5 2V1.5A1.5 1.5 0 0 1 6.5 0h3A1.5 1.5 0 0 1 11 1.5V2h3.5a.5.5 0 0 1 0 1H14l-.84 10.96A1.5 1.5 0 0 1 11.67 15H4.33a1.5 1.5 0 0 1-1.49-1.04L2 3h-.5a.5.5 0 0 1 0-1H5zm1 0h4v-.5a.5.5 0 0 0-.5-.5h-3a.5.5 0 0 0-.5.5V2zM3 3l.82 10.68a.5.5 0 0 0 .5.32h7.34a.5.5 0 0 0 .5-.32L13 3H3z" fill="currentColor"/></svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="!loading && entries.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17 8h1a4 4 0 1 1 0 8h-1"/>
            <path d="M3 8h14v9a4 4 0 0 1-4 4H7a4 4 0 0 1-4-4V8z"/>
            <line x1="6" y1="2" x2="6" y2="4"/>
            <line x1="10" y1="2" x2="10" y2="4"/>
            <line x1="14" y1="2" x2="14" y2="4"/>
          </svg>
        </div>
        <h3>No entries yet</h3>
        <p>Start by adding your first espresso diary entry.</p>
        <RouterLink to="/entry/new" class="btn btn-primary">+ New Entry</RouterLink>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-dot"></div>
        <div class="loading-dot"></div>
        <div class="loading-dot"></div>
      </div>

      <PaginationControls
        :current-page="page"
        :total-pages="totalPages"
        @update:current-page="page = $event"
      />
    </main>

    <ConfirmDialog
      :visible="deleteDialogVisible"
      title="Delete Entry"
      message="Are you sure you want to delete this diary entry? This action cannot be undone."
      @confirm="handleDelete"
      @cancel="deleteDialogVisible = false"
    />
  </div>
</template>

<style scoped>
.page-wrapper {
  min-height: 100vh;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.page-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 28px;
}

.page-title {
  font-size: 28px;
  font-family: var(--font-display);
  font-weight: 400;
  margin: 0;
  letter-spacing: -0.02em;
}

.page-count {
  color: var(--text-dim);
  font-size: 13px;
  margin: 4px 0 0;
  letter-spacing: 0.02em;
}

.page-top-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Filter Bar */
.filter-bar {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 20px 24px;
  margin-bottom: 24px;
}

.filter-row {
  display: flex;
  align-items: flex-end;
  gap: 16px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.filter-group label {
  font-size: 10px;
  font-weight: 600;
  color: var(--text-dim);
  text-transform: uppercase;
  letter-spacing: 0.1em;
}

.filter-group input,
.filter-group select {
  padding: 7px 12px;
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-sm);
  font-family: var(--font-body);
  font-size: 13px;
  color: var(--text);
  background: var(--bg);
  min-width: 120px;
}

.filter-actions {
  display: flex;
  gap: 8px;
  align-items: flex-end;
}

/* Slide transition */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.25s ease;
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  max-height: 0;
  margin-bottom: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.slide-enter-to,
.slide-leave-from {
  max-height: 200px;
}

/* Table */
.entries-table-wrapper {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.entries-table {
  width: 100%;
  border-collapse: collapse;
}

.entries-table th {
  text-align: left;
  padding: 14px 16px;
  font-size: 10px;
  font-weight: 600;
  color: var(--text-dim);
  text-transform: uppercase;
  letter-spacing: 0.1em;
  border-bottom: 1px solid var(--border);
}

.entries-table td {
  padding: 14px 16px;
  font-size: 13px;
  color: var(--text-muted);
  border-bottom: 1px solid var(--border);
}

.entry-row {
  cursor: pointer;
  transition: background var(--transition);
}

.entry-row:hover {
  background: var(--surface-hover);
}

.entry-row:last-child td {
  border-bottom: none;
}

.date-cell {
  display: flex;
  flex-direction: column;
}

.date-main {
  font-weight: 500;
  color: var(--text);
  font-size: 13px;
}

.date-time {
  font-size: 11px;
  color: var(--text-dim);
  margin-top: 2px;
}

.cell-coffee {
  color: var(--accent);
  font-weight: 500;
}

.cell-mono {
  font-variant-numeric: tabular-nums;
  letter-spacing: 0.01em;
}

.actions-cell {
  text-align: right;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-dim);
  padding: 6px;
  border-radius: var(--radius-sm);
  transition: all var(--transition);
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-icon:hover {
  color: var(--error);
  background: var(--error-bg);
}

/* Empty & Loading */
.empty-state {
  text-align: center;
  padding: 80px 24px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
}

.empty-icon {
  color: var(--text-dim);
  margin-bottom: 20px;
}

.empty-state h3 {
  font-family: var(--font-display);
  font-weight: 400;
  font-size: 22px;
  margin: 0 0 8px;
}

.empty-state p {
  color: var(--text-dim);
  margin: 0 0 28px;
  font-size: 14px;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 60px;
}

.loading-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
  animation: pulse 1.2s ease-in-out infinite;
}

.loading-dot:nth-child(2) {
  animation-delay: 0.15s;
}

.loading-dot:nth-child(3) {
  animation-delay: 0.3s;
}

@keyframes pulse {
  0%, 80%, 100% {
    opacity: 0.2;
    transform: scale(0.8);
  }
  40% {
    opacity: 1;
    transform: scale(1);
  }
}

@media (max-width: 900px) {
  .entries-table-wrapper {
    overflow-x: auto;
  }

  .entries-table {
    min-width: 700px;
  }
}
</style>

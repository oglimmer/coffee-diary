<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { diaryService } from '@/services/diaryService'
import { coffeeService } from '@/services/coffeeService'
import { sieveService } from '@/services/sieveService'
import type { Coffee, Sieve, DiaryEntryRequest } from '@/types'
import AppHeader from '@/components/AppHeader.vue'
import StarRating from '@/components/StarRating.vue'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const pageTitle = computed(() => (isEdit.value ? 'Edit Entry' : 'New Entry'))

// Form data
const dateTime = ref('')
const sieveId = ref<number | null>(null)
const temperature = ref(93)
const coffeeId = ref<number | null>(null)
const grindSize = ref<number | null>(null)
const inputWeight = ref<number | null>(null)
const outputWeight = ref<number | null>(null)
const timeSeconds = ref<number | null>(null)
const rating = ref(3)
const notes = ref('')

// Dropdown data
const coffees = ref<Coffee[]>([])
const sieves = ref<Sieve[]>([])

// Inline add
const showNewCoffee = ref(false)
const newCoffeeName = ref('')
const showNewSieve = ref(false)
const newSieveName = ref('')

const error = ref('')
const loading = ref(false)
const saving = ref(false)

function toLocalDateTimeString(date: Date): string {
  const pad = (n: number) => n.toString().padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`
}

async function loadDropdowns() {
  const [c, s] = await Promise.all([coffeeService.getCoffees(), sieveService.getSieves()])
  coffees.value = c
  sieves.value = s
}

async function loadEntry() {
  if (!isEdit.value) {
    dateTime.value = toLocalDateTimeString(new Date())
    return
  }
  loading.value = true
  try {
    const entry = await diaryService.getEntry(Number(route.params.id))
    dateTime.value = toLocalDateTimeString(new Date(entry.dateTime))
    sieveId.value = entry.sieveId ?? null
    temperature.value = entry.temperature
    coffeeId.value = entry.coffeeId ?? null
    grindSize.value = entry.grindSize
    inputWeight.value = entry.inputWeight
    outputWeight.value = entry.outputWeight
    timeSeconds.value = entry.timeSeconds
    rating.value = entry.rating
    notes.value = entry.notes
  } catch {
    error.value = 'Failed to load entry.'
  } finally {
    loading.value = false
  }
}

async function addCoffee() {
  if (!newCoffeeName.value.trim()) return
  try {
    const coffee = await coffeeService.createCoffee(newCoffeeName.value.trim())
    coffees.value.push(coffee)
    coffeeId.value = coffee.id
    newCoffeeName.value = ''
    showNewCoffee.value = false
  } catch {
    error.value = 'Failed to create coffee.'
  }
}

async function addSieve() {
  if (!newSieveName.value.trim()) return
  try {
    const sieve = await sieveService.createSieve(newSieveName.value.trim())
    sieves.value.push(sieve)
    sieveId.value = sieve.id
    newSieveName.value = ''
    showNewSieve.value = false
  } catch {
    error.value = 'Failed to create sieve.'
  }
}

async function handleSubmit() {
  error.value = ''
  saving.value = true

  const request: DiaryEntryRequest = {
    dateTime: dateTime.value,
    sieveId: sieveId.value,
    temperature: temperature.value,
    coffeeId: coffeeId.value,
    grindSize: grindSize.value,
    inputWeight: inputWeight.value,
    outputWeight: outputWeight.value,
    timeSeconds: timeSeconds.value,
    rating: rating.value,
    notes: notes.value,
  }

  try {
    if (isEdit.value) {
      await diaryService.updateEntry(Number(route.params.id), request)
    } else {
      await diaryService.createEntry(request)
    }
    router.push('/')
  } catch (e: unknown) {
    const err = e as { response?: { data?: { message?: string } } }
    error.value = err.response?.data?.message || 'Failed to save entry.'
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadDropdowns()
  loadEntry()
})
</script>

<template>
  <div class="page-wrapper">
    <AppHeader />

    <main class="container fade-up">
      <div class="form-card" v-if="!loading">
        <h2 class="form-title">{{ pageTitle }}</h2>

        <div v-if="error" class="alert alert-error">{{ error }}</div>

        <form @submit.prevent="handleSubmit" class="entry-form">
          <!-- Section: When & Where -->
          <div class="form-section">
            <div class="section-label">Brew Details</div>
            <div class="form-row">
              <div class="form-group">
                <label for="dateTime">Date &amp; Time</label>
                <input id="dateTime" v-model="dateTime" type="datetime-local" required />
              </div>
              <div class="form-group">
                <label for="temperature">Temperature (&deg;C)</label>
                <input id="temperature" v-model.number="temperature" type="number" min="0" max="120" step="1" />
              </div>
            </div>

            <div class="form-row coffee-sieve-row">
              <div class="form-group">
                <label>Coffee</label>
                <div class="inline-add">
                  <select v-model="coffeeId">
                    <option :value="null">-- Select --</option>
                    <option v-for="c in coffees" :key="c.id" :value="c.id">{{ c.name }}</option>
                  </select>
                  <button type="button" class="btn btn-sm btn-outline" @click="showNewCoffee = !showNewCoffee">
                    {{ showNewCoffee ? 'Cancel' : '+ Add' }}
                  </button>
                </div>
                <div v-if="showNewCoffee" class="inline-add-form">
                  <input v-model="newCoffeeName" type="text" placeholder="New coffee name" @keyup.enter="addCoffee" />
                  <button type="button" class="btn btn-sm btn-primary" @click="addCoffee">Add</button>
                </div>
              </div>
              <div class="form-group">
                <label>Sieve</label>
                <div class="inline-add">
                  <select v-model="sieveId">
                    <option :value="null">-- Select --</option>
                    <option v-for="s in sieves" :key="s.id" :value="s.id">{{ s.name }}</option>
                  </select>
                  <button type="button" class="btn btn-sm btn-outline" @click="showNewSieve = !showNewSieve">
                    {{ showNewSieve ? 'Cancel' : '+ Add' }}
                  </button>
                </div>
                <div v-if="showNewSieve" class="inline-add-form">
                  <input v-model="newSieveName" type="text" placeholder="New sieve name" @keyup.enter="addSieve" />
                  <button type="button" class="btn btn-sm btn-primary" @click="addSieve">Add</button>
                </div>
              </div>
            </div>
          </div>

          <!-- Section: Extraction -->
          <div class="form-section">
            <div class="section-label">Extraction</div>
            <div class="form-row four-col">
              <div class="form-group">
                <label for="grindSize">Grind Size</label>
                <input id="grindSize" v-model.number="grindSize" type="number" step="0.1" placeholder="e.g. 2.5" />
              </div>
              <div class="form-group">
                <label for="inputWeight">Input (g)</label>
                <input id="inputWeight" v-model.number="inputWeight" type="number" min="0" step="0.1" placeholder="18.0" />
              </div>
              <div class="form-group">
                <label for="outputWeight">Output (g)</label>
                <input id="outputWeight" v-model.number="outputWeight" type="number" min="0" step="0.1" placeholder="36.0" />
              </div>
              <div class="form-group">
                <label for="timeSeconds">Time (s)</label>
                <input id="timeSeconds" v-model.number="timeSeconds" type="number" min="0" step="1" placeholder="25" />
              </div>
            </div>
          </div>

          <!-- Section: Notes -->
          <div class="form-section">
            <div class="section-label">Evaluation</div>
            <div class="form-group">
              <label>Rating</label>
              <StarRating v-model="rating" :size="32" />
            </div>

            <div class="form-group">
              <label for="notes">Notes</label>
              <textarea id="notes" v-model="notes" rows="4" placeholder="Tasting notes, observations..."></textarea>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-outline" @click="router.push('/')">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ saving ? 'Saving...' : 'Save Entry' }}
            </button>
          </div>
        </form>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="loading-dot"></div>
        <div class="loading-dot"></div>
        <div class="loading-dot"></div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.page-wrapper {
  min-height: 100vh;
}

.container {
  max-width: 720px;
  margin: 0 auto;
  padding: 32px 24px;
}

.form-card {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-xl);
  padding: 40px;
}

.form-title {
  font-size: 28px;
  font-family: var(--font-display);
  font-weight: 400;
  margin: 0 0 32px;
  letter-spacing: -0.02em;
}

.entry-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 24px 0;
  border-bottom: 1px solid var(--border);
}

.form-section:first-child {
  padding-top: 0;
}

.form-section:last-of-type {
  border-bottom: none;
}

.section-label {
  font-size: 10px;
  font-weight: 600;
  color: var(--accent);
  text-transform: uppercase;
  letter-spacing: 0.12em;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-row.coffee-sieve-row {
  grid-template-columns: 3fr 1fr;
}

.form-row.three-col {
  grid-template-columns: 1fr 1fr 1fr;
}

.form-row.four-col {
  grid-template-columns: 1fr 1fr 1fr 1fr;
}

.inline-add {
  display: flex;
  gap: 8px;
}

.inline-add select {
  flex: 1;
}

.inline-add-form {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.inline-add-form input {
  flex: 1;
}

textarea {
  resize: vertical;
  min-height: 80px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
  padding-top: 24px;
  border-top: 1px solid var(--border);
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

@media (max-width: 600px) {
  .form-row,
  .form-row.three-col,
  .form-row.four-col {
    grid-template-columns: 1fr;
  }

  .form-card {
    padding: 24px 20px;
  }
}
</style>

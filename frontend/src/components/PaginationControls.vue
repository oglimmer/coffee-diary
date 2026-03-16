<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPage: number
  totalPages: number
}>()

const emit = defineEmits<{
  'update:currentPage': [page: number]
}>()

const pages = computed(() => {
  const result: (number | string)[] = []
  const total = props.totalPages
  const current = props.currentPage

  if (total <= 7) {
    for (let i = 0; i < total; i++) result.push(i)
    return result
  }

  result.push(0)
  if (current > 2) result.push('...')

  const start = Math.max(1, current - 1)
  const end = Math.min(total - 2, current + 1)
  for (let i = start; i <= end; i++) result.push(i)

  if (current < total - 3) result.push('...')
  result.push(total - 1)

  return result
})

function goTo(page: number) {
  if (page >= 0 && page < props.totalPages) {
    emit('update:currentPage', page)
  }
}
</script>

<template>
  <nav v-if="totalPages > 1" class="pagination">
    <button
      class="page-btn"
      :disabled="currentPage === 0"
      @click="goTo(currentPage - 1)"
    >
      &larr;
    </button>
    <template v-for="(p, idx) in pages" :key="idx">
      <span v-if="p === '...'" class="page-ellipsis">&middot;&middot;&middot;</span>
      <button
        v-else
        class="page-btn"
        :class="{ active: p === currentPage }"
        @click="goTo(p as number)"
      >
        {{ (p as number) + 1 }}
      </button>
    </template>
    <button
      class="page-btn"
      :disabled="currentPage >= totalPages - 1"
      @click="goTo(currentPage + 1)"
    >
      &rarr;
    </button>
  </nav>
</template>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  margin-top: 28px;
}

.page-btn {
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text-muted);
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  font-family: var(--font-body);
  transition: all var(--transition);
  min-width: 36px;
}

.page-btn:hover:not(:disabled):not(.active) {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-glow);
}

.page-btn.active {
  background: var(--accent);
  color: #0f0d0b;
  border-color: var(--accent);
  font-weight: 600;
}

.page-btn:disabled {
  opacity: 0.25;
  cursor: not-allowed;
}

.page-ellipsis {
  color: var(--text-dim);
  padding: 0 4px;
  font-size: 12px;
  letter-spacing: 2px;
}
</style>

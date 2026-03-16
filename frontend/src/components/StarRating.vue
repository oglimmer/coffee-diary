<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: number
    readonly?: boolean
    size?: number
  }>(),
  {
    readonly: false,
    size: 24,
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const ratingLabels: Record<number, string> = {
  1: 'Undrinkable',
  2: 'Ok, but not good',
  3: 'Good',
  4: 'Excellent',
  5: 'Best we ever had',
}

const stars = computed(() => {
  return [1, 2, 3, 4, 5].map((n) => ({
    value: n,
    filled: n <= props.modelValue,
  }))
})

const ratingLabel = computed(() => ratingLabels[props.modelValue] ?? '')

function select(value: number) {
  if (!props.readonly) {
    emit('update:modelValue', value)
  }
}
</script>

<template>
  <div class="star-rating" :class="{ readonly }">
    <button
      v-for="star in stars"
      :key="star.value"
      type="button"
      class="star-btn"
      :class="{ filled: star.filled }"
      :style="{ fontSize: size + 'px', width: size + 6 + 'px', height: size + 6 + 'px' }"
      :disabled="readonly"
      @click="select(star.value)"
    >
      {{ star.filled ? '\u2605' : '\u2606' }}
    </button>
    <span v-if="ratingLabel && !readonly" class="rating-label">{{ ratingLabel }}</span>
  </div>
</template>

<style scoped>
.star-rating {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}

.rating-label {
  margin-left: 10px;
  font-size: 13px;
  color: var(--text-muted);
  font-style: italic;
  font-family: var(--font-display);
}

.star-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--accent);
  padding: 0;
  line-height: 1;
  transition: transform 0.15s ease, color 0.15s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.star-btn:hover:not(:disabled) {
  transform: scale(1.2);
}

.star-btn.filled {
  color: var(--accent);
  text-shadow: 0 0 12px var(--accent-glow-strong);
}

.star-btn:not(.filled) {
  color: var(--text-dim);
}

.readonly .star-btn {
  cursor: default;
}
</style>

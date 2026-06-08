<template>
  <nav v-if="totalPages > 1" class="mt-3">
    <ul class="pagination pagination-sm justify-content-center mb-0">
      <li class="page-item" :class="{ disabled: modelValue <= 1 }">
        <a class="page-link" href="#" @click.prevent="$emit('update:modelValue', modelValue - 1)">上一页</a>
      </li>
      <li v-for="p in pages" :key="p" class="page-item" :class="{ active: p === modelValue }">
        <a class="page-link" href="#" @click.prevent="$emit('update:modelValue', p)">{{ p }}</a>
      </li>
      <li class="page-item" :class="{ disabled: modelValue >= totalPages }">
        <a class="page-link" href="#" @click.prevent="$emit('update:modelValue', modelValue + 1)">下一页</a>
      </li>
    </ul>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: number
  totalPages: number
}>()

defineEmits<{
  'update:modelValue': [value: number]
}>()

const pages = computed(() => {
  const p = props.modelValue
  const tp = props.totalPages
  const delta = 2
  const range: number[] = []
  for (let i = Math.max(2, p - delta); i <= Math.min(tp - 1, p + delta); i++) {
    range.push(i)
  }
  const res: number[] = [1]
  if (range[0] > 2) res.push(-1)
  res.push(...range)
  if (range[range.length - 1] < tp - 1) res.push(-1)
  if (tp > 1) res.push(tp)
  return res
})
</script>

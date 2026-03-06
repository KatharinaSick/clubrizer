<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  title: string
  message: string
  variant?: 'error' | 'warning'
}>()

const variantClass = computed(() => {
  if (props.variant === 'warning') {
    return 'alertWarning'
  }
  return 'alertError'
})
</script>

<template>
  <div class="alert" :class="variantClass">
    <p style="font-weight: var(--font-weight-bold)">{{title}}</p>
    <p>{{message}}</p>
  </div>
</template>

<style scoped>
.alert {
  width: 100%;
  padding: var(--padding);
  border-radius: var(--border-radius);
  box-sizing: border-box;
}

.alertError {
  background: var(--light-red);
  border: 1px solid var(--red);
}

.alertWarning {
  background: var(--light-orange);
  border: 1px solid var(--orange);
}
</style>

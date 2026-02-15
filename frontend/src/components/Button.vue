<script setup lang="ts">
defineProps<{
  title: string
  loading?: boolean
  disabled?: boolean
}>()
</script>

<template>
  <button class="button" :disabled="disabled || loading">
    <span v-if="loading">Loading...</span>
    <span v-else>{{ title }}</span>
  </button>
</template>

<style scoped>
.button {
  width: 100%;
  padding: 12px;
  margin-top: 24px;

  background: var(--horizotal-gradient);
  border-radius: var(--border-radius);
  border: 0;

  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-medium);
  color: var(--white);
  cursor: pointer;
}

.button:active {
  background: var(--horizotal-gradient-active);
}

.button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background: var(--horizotal-gradient-active);
}
</style>

<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <button class="button">{{ title }}</button>
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
}

.button:active {
  background: var(--horizotal-gradient-active);
}
</style>

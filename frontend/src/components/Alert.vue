<script setup lang="ts">
defineProps<{
  title: string
  message: string
}>()
</script>

<template>
  <div class="alert">
    <strong>{{title}}</strong>
    <p>{{message}}</p>
  </div>
</template>

<style scoped>
.alert {
  width: 100%;
  background: var(--light-red);
  padding: var(--padding);
  border: 1px solid var(--red);
  border-radius: var(--border-radius);
  box-sizing: border-box;
}
</style>

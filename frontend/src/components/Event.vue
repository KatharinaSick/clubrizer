<script setup lang="ts">

interface Category {
  name: string
  color: string
  picture?: string
}

export interface EventProps {
  title: string
  description: string
  category: Category
  location: string
  startTime: Date
}

defineProps<{
  event: EventProps
}>()
</script>

<template>
  <div class="event">
    {{event.title}}
  </div>
</template>

<style scoped>
.event {
  width: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: var(--padding);
  background-color: var(--white);
  margin: 0;
}
</style>

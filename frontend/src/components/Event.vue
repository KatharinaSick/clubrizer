<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import EventTitle from '@/components/EventTitle.vue'

interface Category {
  name: string
  color: string
  picture?: string
}

export interface EventProps {
  id: string
  title: string
  description: string
  category: Category
  location: string
  startTime: string | Date
}

const props = defineProps<{
  event: EventProps
}>()

const emit = defineEmits(['click'])

const handleClick = (event: MouseEvent) => {
  emit('click', event)
}
</script>

<template>
  <div class="event" @click="handleClick">
    <img class="image" :src="event.category.picture" />
    <EventTitle class="details" :event="event" />
  </div>
</template>

<style scoped>
.event {
  width: 100%;
  margin-bottom: var(--gap);

  box-shadow: var(--box-shadow);
  border-radius: var(--border-radius);
}

.image {
  width: 100%;
  height: 128px;

  object-fit: cover;

  margin: 0;
  padding: 0;

  display: block;

  border-top-left-radius: var(--border-radius);
  border-top-right-radius: var(--border-radius);
}

.details {
  width: 100%;
  margin: 0;
  padding: 0;

  display: flex;
  flex-direction: row;

  background-color: var(--background-color);
  border-bottom-left-radius: var(--border-radius);
  border-bottom-right-radius: var(--border-radius);
  overflow: hidden;
}
</style>

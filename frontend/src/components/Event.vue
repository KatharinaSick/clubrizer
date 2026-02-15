<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

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
  startTime: string | Date
}

const props = defineProps<{
  event: EventProps
}>()

const { locale } = useI18n()

const parsedDate = computed(() => {
  const date = new Date(props.event.startTime)

  // Use Intl.DateTimeFormat for consistent localization
  const monthFormatter = new Intl.DateTimeFormat(locale.value, { month: 'short' })
  const dayFormatter = new Intl.DateTimeFormat(locale.value, { day: 'numeric' })
  const timeFormatter = new Intl.DateTimeFormat(locale.value, {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false // Force 24h format or remove to respect locale default
  })

  return {
    month: monthFormatter.format(date),
    day: dayFormatter.format(date),
    time: timeFormatter.format(date)
  }
})
</script>

<template>
  <div class="event">
    <img class="image" :src="event.category.picture" />
    <div class="details">
      <div class="date" :style="{ background: event.category.color }">
        <span class="month">{{ parsedDate.month }}</span>
        <span class="day">{{ parsedDate.day }}</span>
      </div>
      <div class="text">
        <span class="name">{{ event.title }}</span>
        <span class="time">{{ parsedDate.time }} | {{ event.location }}</span>
      </div>
    </div>
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

  border-bottom-left-radius: var(--border-radius);
  border-bottom-right-radius: var(--border-radius);
}

.date {
  width: 36px;
  height: 36px;
  padding: var(--padding);

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  border-bottom-left-radius: var(--border-radius);
}

.month {
  color: var(--background-color);
  font-size: var(--font-size-small);
  text-transform: capitalize;
}

.day {
  color: var(--background-color);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
}

.text {
  width: 100%;
  padding-left: var(--padding);
  display: flex;
  flex-direction: column;
  justify-content: center;

  border-bottom-left-radius: var(--border-radius);
}

.name {
  font-weight: var(--font-weight-bold);
}
</style>

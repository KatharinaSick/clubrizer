<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

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
  <div class="eventTitle">
    <div class="eventTitleDate" :style="{ background: event.category.color }">
      <span class="eventTitleMonth">{{ parsedDate.month }}</span>
      <span class="eventTitleDay">{{ parsedDate.day }}</span>
    </div>
    <div class="eventTitleText">
      <span class="eventTitleName">{{ event.title }}</span>
      <span>{{ parsedDate.time }} | {{ event.location }}</span>
    </div>
  </div>
</template>

<style scoped>
.eventTitle {
  width: 100%;
  margin: 0;
  padding: 0;

  display: flex;
  flex-direction: row;
}

.eventTitleDate {
  width: 36px;
  height: 36px;
  padding: var(--padding);

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.eventTitleMonth {
  color: var(--background-color);
  font-size: var(--font-size-small);
  text-transform: capitalize;
}

.eventTitleDay {
  color: var(--background-color);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
}

.eventTitleText {
  width: 100%;
  padding-left: var(--padding);
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.eventTitleName {
  font-weight: var(--font-weight-bold);
}
</style>

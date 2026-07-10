<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import EventTitle from '@/components/EventTitle.vue'

interface Category {
  name: string
  color: string
  picture?: string
  customLabel?: string
}

interface Creator {
  givenName: string
  familyName: string
  nickName: string
}

export interface EventProps {
  id: string
  title: string
  description: string
  category: Category
  location: string
  startTime: string | Date
  creator: Creator
  cancelledAt?: string
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
    <div class="eventImageContainer">
      <img class="eventImage" :class="{ eventImageCancelled: event.cancelledAt }" :src="event.category.picture" />
      <span v-if="event.cancelledAt" class="eventImageBadge eventImageBadgeCancelled">{{ $t('events.cancelled') }}</span>
      <span v-else-if="event.category.customLabel" class="eventImageBadge">{{ event.category.customLabel }}</span>
    </div>
    <EventTitle class="details" :class="{ eventDetailsCancelled: event.cancelledAt }" :event="event" />
  </div>
</template>

<style scoped>
.event {
  width: 100%;
  margin-bottom: var(--gap);

  box-shadow: var(--box-shadow);
  border-radius: var(--border-radius);
}

.eventImageContainer {
  position: relative;
}

.eventImage {
  width: 100%;
  height: 128px;

  object-fit: cover;

  margin: 0;
  padding: 0;

  display: block;

  border-top-left-radius: var(--border-radius);
  border-top-right-radius: var(--border-radius);
}

.eventImageCancelled {
  filter: grayscale(1) opacity(0.6);
}

.eventDetailsCancelled {
  opacity: 0.6;
}

.eventImageBadge {
  position: absolute;
  bottom: var(--gap);
  left: var(--gap);
  padding: 2px var(--gap);
  border-radius: var(--border-radius);
  background-color: var(--white);
  color: var(--text-color);
  font-size: var(--font-size-small);
}

.eventImageBadgeCancelled {
  background-color: var(--red);
  color: var(--white);
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

.details :deep(.eventTitleName) {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>

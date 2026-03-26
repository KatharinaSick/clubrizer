<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed, onMounted, ref } from 'vue'
import axios from '@/plugins/axios'
import Button from '@/components/Button.vue'
import Header from '@/components/Header.vue'
import EventTitle from '@/components/EventTitle.vue'
import type { Event } from '@/service/events'
import i18n from '@/plugins/i18n'
import IconError from '@/components/icons/IconError.vue'
import IconCheckmark from '@/components/icons/IconCheckmark.vue'
import IconMapMarker from '@/components/icons/IconMapMarker.vue'
import Divider from '@/components/Divider.vue'
import RequestError from '@/components/RequestError.vue'

const route = useRoute()
const eventId = route.params.id as string

const event = ref<Event | null>(null)

onMounted(() => {
  loadEvent()
})

const loadEvent = () => {
  axios.get(`/events/${eventId}`)
    .then(response => {
      event.value = response.data
    })
}

const formattedStartTime = computed(() => {
  if (!event.value?.startTime) return ''
  const date = new Date(event.value.startTime)
  const datePart = date.toLocaleDateString(navigator.language, {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
  })
  const timePart = date.toLocaleTimeString(navigator.language, {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
  return `${datePart} ${i18n.global.t('events.detail.at')} ${timePart}`
})
</script>

<template>
  <div>
    <Header
      left-action="back"
    />

    <RequestError />

    <div v-if="event">
      <!-- Hero Image -->
      <img :src="event.category.picture" class="eventDetailImage" />

      <div class="eventDetailContent">
        <!-- Title & Category -->
        <EventTitle class="eventDetailTitle" :event="event" />

        <div class="eventDetailInfo">
          <!-- Attendance Buttons -->
          <div class="eventDetailAttendanceButtons">
            <Button :title="i18n.global.t('events.detail.wontGo')" theme="secondary" disabled>
              <template #icon>
                <IconError style="color: var(--red)" />
              </template>
            </Button>
            <Button :title="i18n.global.t('events.detail.going')" theme="secondary" disabled>
              <template #icon>
                <IconCheckmark style="color: var(--green)" />
              </template>
            </Button>
          </div>

          <!-- Date & Location -->
          <p class="eventDetailDate">{{ formattedStartTime }}</p>
          <div class="eventDetailLocation">
            <IconMapMarker class="eventDetailLocationIcon" />
            <span>{{ event.location }}</span>
          </div>

          <!-- Description -->
          <template v-if="event.description">
            <Divider class="eventDetailDivider" />
            <p class="eventDetailDescription">{{ event.description }}</p>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.eventDetailImage {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 180px;
  object-fit: cover;
  z-index: -1;
}

.eventDetailContent {
  position: absolute;
  width: 100%;
  top: 180px;
  left: 0;
}

.eventDetailTitle {
  width: 100vw;
  box-shadow: var(--box-shadow);
}

.eventDetailInfo {
  width: 100%;
  padding: var(--padding);
  box-sizing: border-box;
}

.eventDetailAttendanceButtons {
  width: 100%;
  display: flex;
  flex-direction: row;
  gap: var(--gap);
  margin: var(--padding) 0 var(--padding) 0;
}

.eventDetailDate {
  padding-top: var(--padding);
  font-weight: var(--font-weight-medium);
}

.eventDetailLocation {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: var(--gap);
  padding-top: var(--gap);
}

.eventDetailLocationIcon {
  flex-shrink: 0;
  color: var(--blue);
}

.eventDetailDivider {
  margin: var(--padding) 0;
}

</style>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed, onMounted, ref } from 'vue'
import axios from '@/plugins/axios'
import Avatar from '@/components/Avatar.vue'
import Button from '@/components/Button.vue'
import Header from '@/components/Header.vue'
import EventTitle from '@/components/EventTitle.vue'
import type { EventDetail } from '@/service/events'
import { upsertEventResponse } from '@/service/events'
import i18n from '@/plugins/i18n'
import IconError from '@/components/icons/IconError.vue'
import IconCheckmark from '@/components/icons/IconCheckmark.vue'
import IconMapMarker from '@/components/icons/IconMapMarker.vue'
import Divider from '@/components/Divider.vue'
import RequestError from '@/components/RequestError.vue'

const route = useRoute()
const eventId = route.params.id as string

const event = ref<EventDetail | null>(null)
const pendingResponse = ref<boolean | null>(null)

onMounted(() => {
  loadEvent()
})

const loadEvent = async () => {
  const response = await axios.get(`/events/${eventId}`)
  event.value = response.data
}

const isPastEvent = computed(() => {
  if (!event.value?.startTime) return false
  return new Date(event.value.startTime) <= new Date()
})

const submitResponse = async (response: boolean) => {
  if (!event.value || isPastEvent.value || event.value.responses?.currentUserResponse === response || pendingResponse.value !== null) return
  pendingResponse.value = response
  try {
    await upsertEventResponse(eventId, response)
    // Reload the full event to keep counts and attendee list in sync with the server
    await loadEvent()
  } finally {
    pendingResponse.value = null
  }
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
            <Button
              :title="i18n.global.t('events.detail.wontGo')"
              :theme="event.responses?.currentUserResponse === false ? 'red' : 'secondary'"
              :loading="pendingResponse === false"
              :disabled="isPastEvent || pendingResponse !== null"
              @click="submitResponse(false)"
            >
              <template #icon>
                <IconError :class="{ 'eventDetailIconNotGoing': event.responses?.currentUserResponse !== false }" />
              </template>
            </Button>
            <Button
              :title="i18n.global.t('events.detail.going')"
              :theme="event.responses?.currentUserResponse === true ? 'green' : 'secondary'"
              :loading="pendingResponse === true"
              :disabled="isPastEvent || pendingResponse !== null"
              @click="submitResponse(true)"
            >
              <template #icon>
                <IconCheckmark :class="{ 'eventDetailIconGoing': event.responses?.currentUserResponse !== true }" />
              </template>
            </Button>
          </div>

          <!-- Date & Location -->
          <p class="eventDetailDate">{{ formattedStartTime }}</p>
          <div class="eventDetailLocation">
            <IconMapMarker class="eventDetailLocationIcon" />
            <span>{{ event.location }}</span>
          </div>

          <!-- Attendees -->
          <Divider class="eventDetailDivider" />
          <template v-if="event.responses && (event.responses.going + event.responses.notGoing) > 0">
            <p class="eventDetailAttendeeCounts">
              {{ i18n.global.t('events.detail.attendees.going', { count: event.responses.going }) }}
              &middot;
              {{ i18n.global.t('events.detail.attendees.notGoing', { count: event.responses.notGoing }) }}
            </p>
            <div class="eventDetailAvatarGrid">
              <div
                v-for="attendee in event.responses.attendees"
                :key="attendee.id"
                class="eventDetailAvatarWrapper"
              >
                <Avatar
                  :picture="attendee.picture"
                  :given-name="attendee.givenName"
                  :family-name="attendee.familyName"
                  :nick-name="attendee.nickName"
                  :label="attendee.nickName || attendee.givenName"
                  size="md"
                  :class="{ 'eventDetailAvatarNotGoing': !attendee.response }"
                />
                <span
                  class="eventDetailAvatarBadge"
                  :class="attendee.response ? 'eventDetailAvatarBadgeGoing' : 'eventDetailAvatarBadgeNotGoing'"
                />
              </div>
            </div>
          </template>
          <p v-else class="eventDetailNoResponses">
            {{ i18n.global.t('events.detail.attendees.noResponses') }}
          </p>

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
  margin: var(--padding) 0;
}

.eventDetailIconGoing {
  color: var(--green);
}

.eventDetailIconNotGoing {
  color: var(--red);
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

.eventDetailAttendeeCounts {
  font-weight: var(--font-weight-medium);
  padding-bottom: var(--gap);
}

.eventDetailNoResponses {
  color: var(--text-gray);
}

.eventDetailAvatarGrid {
  display: flex;
  flex-wrap: wrap;
  gap: var(--padding);
}

.eventDetailAvatarWrapper {
  position: relative;
}

.eventDetailAvatarNotGoing {
  opacity: 0.4;
}

.eventDetailAvatarBadge {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: 2px solid var(--background-color);
}

.eventDetailAvatarBadgeGoing {
  background: var(--green);
}

.eventDetailAvatarBadgeNotGoing {
  background: var(--red);
}

</style>

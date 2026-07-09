<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { computed, onMounted, ref } from 'vue'
import axios from '@/plugins/axios'
import Avatar from '@/components/Avatar.vue'
import Button from '@/components/Button.vue'
import EventTitle from '@/components/EventTitle.vue'
import type { EventDetail } from '@/service/events'
import { upsertEventResponse, deleteEvent } from '@/service/events'
import i18n from '@/plugins/i18n'
import IconBack from '@/components/icons/IconBack.vue'
import IconError from '@/components/icons/IconError.vue'
import IconCheckmark from '@/components/icons/IconCheckMark.vue'
import IconMapMarker from '@/components/icons/IconMapMarker.vue'
import Divider from '@/components/Divider.vue'
import RequestError from '@/components/RequestError.vue'
import MenuButton from '@/components/MenuButton.vue'
import type { MenuItem } from '@/components/MenuButton.vue'
import Modal from '@/components/Modal.vue'

const route = useRoute()
const router = useRouter()
const eventId = route.params.id as string

const event = ref<EventDetail | null>(null)
const pendingResponse = ref<boolean | null>(null)
const showDeleteConfirm = ref(false)
const isDeleting = ref(false)

const menuItems = computed<MenuItem[]>(() => [
  {
    label: i18n.global.t('events.detail.menu.delete'),
    danger: true,
    onClick: () => { showDeleteConfirm.value = true },
  },
])

const confirmDelete = async () => {
  if (isDeleting.value) return
  isDeleting.value = true
  try {
    await deleteEvent(eventId)
    router.replace('/events')
  } catch {
    // Reveal the global error banner behind the modal
    showDeleteConfirm.value = false
  } finally {
    isDeleting.value = false
  }
}

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
  <div class="eventDetailPage">
    <div class="eventDetailScroll">
      <button class="eventDetailBackButton" @click="router.back()">
        <IconBack class="eventDetailBackIcon" />
      </button>

      <div v-if="event?.canDelete" class="eventDetailMenuButton">
        <MenuButton :items="menuItems" :aria-label="i18n.global.t('events.detail.menu.label')" />
      </div>

      <RequestError class="eventDetailRequestError" />

      <div v-if="event">
        <!-- Hero Image -->
        <div class="eventDetailHero">
          <img :src="event.category.picture" class="eventDetailImage" />
          <span v-if="event.category.customLabel" class="eventDetailBadge">{{ event.category.customLabel }}</span>
        </div>

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
          <Divider class="eventDetailDivider" />
          <div class="eventDetailCreator">
            <Avatar
              :picture="event.creator.picture"
              :given-name="event.creator.givenName"
              :family-name="event.creator.familyName"
              :nick-name="event.creator.nickName"
              size="sm"
            />
            <span>{{ i18n.global.t('events.createdBy') }} {{ event.creator.nickName || event.creator.givenName }}</span>
          </div>

          <!-- Description -->
          <p v-if="event.description" class="eventDetailDescription">{{ event.description }}</p>

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
        </div>
      </div>
    </div>

    <Modal v-if="showDeleteConfirm">
      <div class="eventDetailDeleteConfirm">
        <h2 class="eventDetailDeleteTitle">{{ i18n.global.t('events.detail.deleteConfirm.title') }}</h2>
        <p class="eventDetailDeleteMessage">{{ i18n.global.t('events.detail.deleteConfirm.message') }}</p>
        <div class="eventDetailDeleteActions">
          <Button
            :title="i18n.global.t('events.detail.deleteConfirm.confirm')"
            theme="red"
            :loading="isDeleting"
            @click="confirmDelete"
          />
          <Button
            :title="i18n.global.t('events.detail.deleteConfirm.cancel')"
            theme="secondary"
            :disabled="isDeleting"
            @click="showDeleteConfirm = false"
          />
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.eventDetailPage {
  position: relative;
  height: 100%;
}

.eventDetailScroll {
  position: absolute;
  inset: 0;
  overflow-y: auto;
}

.eventDetailBackButton {
  position: absolute;
  top: var(--padding);
  left: var(--padding);
  z-index: 10;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(4px);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--box-shadow);
  padding: 0;
}

.eventDetailBackIcon {
  width: 22px;
  height: 22px;
  color: var(--text-color);
}

.eventDetailMenuButton {
  position: absolute;
  top: var(--padding);
  right: var(--padding);
  z-index: 10;
}

.eventDetailDeleteConfirm {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.eventDetailDeleteTitle {
  margin: 0;
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.eventDetailDeleteMessage {
  margin: 0;
  line-height: 1.5;
}

.eventDetailDeleteActions {
  margin-top: var(--padding);
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.eventDetailRequestError {
  padding: var(--padding);
}

.eventDetailHero {
  position: relative;
  width: 100%;
  height: 180px;
}

.eventDetailImage {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.eventDetailBadge {
  position: absolute;
  bottom: var(--gap);
  left: var(--padding);
  padding: 2px var(--gap);
  border-radius: var(--border-radius);
  background-color: var(--white);
  color: var(--text-color);
  font-size: var(--font-size-small);
}

.eventDetailTitle {
  width: 100%;
  box-shadow: var(--box-shadow);
}

.eventDetailTitle :deep(.eventTitleName) {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.eventDetailInfo {
  width: 100%;
  padding: var(--padding);
  box-sizing: border-box;
}

.eventDetailCreator {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: var(--gap);
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

.eventDetailDescription {
  margin-top: var(--gap);
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

.eventDetailAvatarNotGoing :deep(.avatarImage),
.eventDetailAvatarNotGoing :deep(.avatarFallback) {
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

@media (min-width: 768px) {
  .eventDetailPage {
    height: auto;
    padding: var(--padding);
    box-sizing: border-box;
  }

  .eventDetailScroll {
    position: relative;
    overflow-y: visible;
    max-width: var(--content-max-width);
    margin: 0 auto;
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
  }

  .eventDetailHero {
    height: 320px;
  }

  .eventDetailImage {
    border-top-left-radius: var(--border-radius);
    border-top-right-radius: var(--border-radius);
  }

  .eventDetailInfo {
    border-bottom-left-radius: var(--border-radius);
    border-bottom-right-radius: var(--border-radius);
  }
}
</style>

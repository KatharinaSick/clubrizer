<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { computed, onMounted, ref } from 'vue'
import { useMediaQuery } from '@vueuse/core'
import axios from '@/plugins/axios'
import Alert from '@/components/Alert.vue'
import Avatar from '@/components/Avatar.vue'
import Button from '@/components/Button.vue'
import EventTitle from '@/components/EventTitle.vue'
import type { EventDetail, MyResponse } from '@/service/events'
import { upsertEventResponse, deleteEvent, cancelEvent, uncancelEvent } from '@/service/events'
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
import UserProfileModal from '@/components/UserProfileModal.vue'

type UserForModal = {
  givenName: string
  familyName: string
  nickName: string | null
  picture?: string | null
  isKid?: boolean
  parent?: string
}

const route = useRoute()
const router = useRouter()
const eventId = route.params.id as string

const event = ref<EventDetail | null>(null)
const pendingResponse = ref<boolean | null>(null)
const submittingIds = ref<Set<string>>(new Set())
const showDeleteConfirm = ref(false)
const isDeleting = ref(false)
const showCancelConfirm = ref(false)
const isCancelling = ref(false)
const showRestoreConfirm = ref(false)
const isRestoring = ref(false)
const selectedUser = ref<UserForModal | null>(null)

const isCancelled = computed(() => !!event.value?.cancelledAt)

// Below the layout breakpoint the per-person RSVP toggles show icon-only to stay compact;
// above it they spell out the action.
const isCompact = useMediaQuery('(max-width: 767px)')

const menuItems = computed<MenuItem[]>(() => {
  const items: MenuItem[] = []
  if (event.value?.canCancel && !isCancelled.value) {
    items.push({
      label: i18n.global.t('events.detail.menu.cancel'),
      onClick: () => { showCancelConfirm.value = true },
    })
  }
  if (event.value?.canCancel && isCancelled.value) {
    items.push({
      label: i18n.global.t('events.detail.menu.restore'),
      onClick: () => { showRestoreConfirm.value = true },
    })
  }
  if (event.value?.canDelete) {
    items.push({
      label: i18n.global.t('events.detail.menu.delete'),
      danger: true,
      onClick: () => { showDeleteConfirm.value = true },
    })
  }
  return items
})

const confirmRestore = async () => {
  if (isRestoring.value) return
  isRestoring.value = true
  try {
    await uncancelEvent(eventId)
    await loadEvent()
    showRestoreConfirm.value = false
  } catch {
    showRestoreConfirm.value = false
  } finally {
    isRestoring.value = false
  }
}

const confirmCancel = async () => {
  if (isCancelling.value) return
  isCancelling.value = true
  try {
    await cancelEvent(eventId)
    await loadEvent()
    showCancelConfirm.value = false
  } catch {
    showCancelConfirm.value = false
  } finally {
    isCancelling.value = false
  }
}

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

// The people the current account can RSVP for on this event: the account holder (if
// they participate) plus each approved kid, each with their current response.
const myPeople = computed(() => event.value?.responses?.myResponses ?? [])
const isSinglePerson = computed(() => myPeople.value.length === 1)
const soloPerson = computed(() => myPeople.value[0] ?? null)

// What a Going/Not-going button should reflect: for a single person it mirrors their
// current response; for multiple people the buttons just open the picker.
const soloGoing = computed(() => isSinglePerson.value && soloPerson.value?.response === true)
const soloNotGoing = computed(() => isSinglePerson.value && soloPerson.value?.response === false)

const submitSolo = async (response: boolean) => {
  const person = soloPerson.value
  if (!person || person.response === response || pendingResponse.value !== null) return
  pendingResponse.value = response
  try {
    await upsertEventResponse(eventId, response, person.isSelf ? undefined : person.id)
    await loadEvent()
  } finally {
    pendingResponse.value = null
  }
}

// Multi-person accounts: set one person's response directly from their inline row. Each
// row tracks its own in-flight state so several can be tapped without blocking each other.
const setPersonResponse = async (person: MyResponse, response: boolean) => {
  if (isPastEvent.value || person.response === response || submittingIds.value.has(person.id)) return
  submittingIds.value = new Set(submittingIds.value).add(person.id)
  try {
    await upsertEventResponse(eventId, response, person.isSelf ? undefined : person.id)
    await loadEvent()
  } finally {
    const next = new Set(submittingIds.value)
    next.delete(person.id)
    submittingIds.value = next
  }
}

const sortedAttendees = computed(() => {
  const attendees = event.value?.responses?.attendees ?? []
  return [...attendees].sort((a, b) => (b.response ? 1 : 0) - (a.response ? 1 : 0))
})

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

      <div v-if="menuItems.length > 0" class="eventDetailMenuButton">
        <MenuButton :items="menuItems" :aria-label="i18n.global.t('events.detail.menu.label')" />
      </div>

      <RequestError class="eventDetailRequestError" />

      <div v-if="event">
        <!-- Hero Image -->
        <div class="eventDetailHero">
          <img :src="event.category.picture" class="eventDetailImage" :class="{ eventDetailImageCancelled: isCancelled }" />
          <span v-if="event.category.customLabel" class="eventDetailBadge">{{ event.category.customLabel }}</span>
        </div>

        <!-- Title & Category -->
        <EventTitle class="eventDetailTitle" :class="{ eventDetailTitleCancelled: isCancelled }" :event="event" />

        <div class="eventDetailInfo">
          <!-- Cancelled banner -->
          <Alert
            v-if="isCancelled"
            :message="i18n.global.t('events.detail.cancelledBanner')"
            variant="error"
            class="eventDetailCancelledBanner"
          />

          <div :class="{ eventDetailInfoBodyCancelled: isCancelled }">
          <!-- RSVP controls (hidden when the account has no one to RSVP, e.g. a guardian
               with no approved kids yet). Single person → prominent buttons; multiple
               people → a compact row per person so the whole family is visible and settable
               at a glance, each with its own going/not-going toggle. -->
          <template v-if="!isCancelled && myPeople.length > 0">
            <div v-if="isSinglePerson" class="eventDetailAttendanceButtons">
              <Button
                :title="i18n.global.t('events.detail.wontGo')"
                :theme="soloNotGoing ? 'red' : 'secondary'"
                :loading="pendingResponse === false"
                :disabled="isPastEvent || pendingResponse !== null"
                @click="submitSolo(false)"
              >
                <template #icon>
                  <IconError :class="{ 'eventDetailIconNotGoing': !soloNotGoing }" />
                </template>
              </Button>
              <Button
                :title="i18n.global.t('events.detail.going')"
                :theme="soloGoing ? 'green' : 'secondary'"
                :loading="pendingResponse === true"
                :disabled="isPastEvent || pendingResponse !== null"
                @click="submitSolo(true)"
              >
                <template #icon>
                  <IconCheckmark :class="{ 'eventDetailIconGoing': !soloGoing }" />
                </template>
              </Button>
            </div>
            <div v-else class="eventDetailPeople">
              <div v-for="person in myPeople" :key="person.id" class="eventDetailPerson">
                <Avatar :picture="person.picture" :given-name="person.name" :family-name="null" size="sm" />
                <span class="eventDetailPersonName">{{ person.name }}</span>
                <div class="eventDetailPersonChoices">
                  <Button
                    inline
                    :icon-only="isCompact"
                    :theme="person.response === false ? 'red' : 'secondary'"
                    :title="i18n.global.t('events.detail.wontGo')"
                    :disabled="isPastEvent || submittingIds.has(person.id)"
                    @click="setPersonResponse(person, false)"
                  >
                    <template #icon>
                      <IconError :class="{ eventDetailIconNotGoing: person.response !== false }" />
                    </template>
                  </Button>
                  <Button
                    inline
                    :icon-only="isCompact"
                    :theme="person.response === true ? 'green' : 'secondary'"
                    :title="i18n.global.t('events.detail.going')"
                    :disabled="isPastEvent || submittingIds.has(person.id)"
                    @click="setPersonResponse(person, true)"
                  >
                    <template #icon>
                      <IconCheckmark :class="{ eventDetailIconGoing: person.response !== true }" />
                    </template>
                  </Button>
                </div>
              </div>
            </div>
          </template>

          <!-- Date & Location -->
          <p class="eventDetailDate">{{ formattedStartTime }}</p>
          <div class="eventDetailLocation">
            <IconMapMarker class="eventDetailLocationIcon" />
            <span>{{ event.location }}</span>
          </div>
          <Divider class="eventDetailDivider" />
          <div class="eventDetailCreator">
            <button class="eventDetailAvatarButton" @click="selectedUser = event.creator">
              <Avatar
                :picture="event.creator.picture"
                :given-name="event.creator.givenName"
                :family-name="event.creator.familyName"
                :nick-name="event.creator.nickName"
                size="sm"
              />
            </button>
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
              <button
                v-for="attendee in sortedAttendees"
                :key="attendee.id"
                class="eventDetailAvatarWrapper"
                @click="selectedUser = attendee"
              >
                <Avatar
                  :picture="attendee.picture"
                  :given-name="attendee.givenName"
                  :family-name="attendee.familyName"
                  :nick-name="attendee.nickName"
                  size="md"
                  :class="{ 'eventDetailAvatarGoing': attendee.response, 'eventDetailAvatarNotGoing': !attendee.response }"
                />
                <span
                  class="eventDetailAvatarBadge"
                  :class="attendee.response ? 'eventDetailAvatarBadgeGoing' : 'eventDetailAvatarBadgeNotGoing'"
                />
              </button>
            </div>
          </template>
          <p v-else class="eventDetailNoResponses">
            {{ i18n.global.t('events.detail.attendees.noResponses') }}
          </p>
          </div>
        </div>
      </div>
    </div>

    <UserProfileModal
      v-if="selectedUser"
      :given-name="selectedUser.givenName"
      :family-name="selectedUser.familyName"
      :nick-name="selectedUser.nickName"
      :picture="selectedUser.picture"
      :parent="selectedUser.isKid ? selectedUser.parent : undefined"
      @close="selectedUser = null"
    />

    <Modal v-if="showCancelConfirm">
      <div class="eventDetailCancelConfirm">
        <h2 class="eventDetailDeleteTitle">{{ i18n.global.t('events.detail.cancelConfirm.title') }}</h2>
        <p class="eventDetailDeleteMessage">{{ i18n.global.t('events.detail.cancelConfirm.message') }}</p>
        <div class="eventDetailDeleteActions">
          <Button
            :title="i18n.global.t('events.detail.cancelConfirm.confirm')"
            theme="red"
            :loading="isCancelling"
            @click="confirmCancel"
          />
          <Button
            :title="i18n.global.t('events.detail.cancelConfirm.back')"
            theme="secondary"
            :disabled="isCancelling"
            @click="showCancelConfirm = false"
          />
        </div>
      </div>
    </Modal>

    <Modal v-if="showRestoreConfirm">
      <div class="eventDetailRestoreConfirm">
        <h2 class="eventDetailDeleteTitle">{{ i18n.global.t('events.detail.restoreConfirm.title') }}</h2>
        <p class="eventDetailDeleteMessage">{{ i18n.global.t('events.detail.restoreConfirm.message') }}</p>
        <div class="eventDetailDeleteActions">
          <Button
            :title="i18n.global.t('events.detail.restoreConfirm.confirm')"
            theme="primary"
            :loading="isRestoring"
            @click="confirmRestore"
          />
          <Button
            :title="i18n.global.t('events.detail.restoreConfirm.back')"
            theme="secondary"
            :disabled="isRestoring"
            @click="showRestoreConfirm = false"
          />
        </div>
      </div>
    </Modal>

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

.eventDetailDeleteConfirm,
.eventDetailRestoreConfirm {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.eventDetailImageCancelled {
  filter: grayscale(1) opacity(0.5);
}

.eventDetailTitleCancelled {
  opacity: 0.5;
}

.eventDetailInfoBodyCancelled {
  opacity: 0.5;
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

.eventDetailPeople {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  margin: var(--padding) 0;
}

.eventDetailPerson {
  display: flex;
  align-items: center;
  gap: var(--gap);
}

.eventDetailPersonName {
  flex: 1;
  min-width: 0;
  font-weight: var(--font-weight-medium);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.eventDetailPersonChoices {
  display: flex;
  gap: var(--gap);
  flex-shrink: 0;
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

.eventDetailCancelledBanner {
  margin-top: var(--padding);
}

.eventDetailCancelConfirm {
  display: flex;
  flex-direction: column;
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

.eventDetailAvatarButton {
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  border-radius: 50%;
  display: flex;
}

.eventDetailAvatarWrapper {
  position: relative;
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  border-radius: 50%;
}

.eventDetailAvatarGoing :deep(.avatarImage),
.eventDetailAvatarGoing :deep(.avatarFallback) {
  box-shadow: 0 0 0 2px var(--gray);
}

.eventDetailAvatarNotGoing :deep(.avatarImage),
.eventDetailAvatarNotGoing :deep(.avatarFallback) {
  filter: grayscale(1) opacity(0.5);
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

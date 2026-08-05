<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { computed, onMounted, ref } from 'vue'
import { useMediaQuery } from '@vueuse/core'
import axios from '@/plugins/axios'
import Alert from '@/components/Alert.vue'
import Avatar from '@/components/Avatar.vue'
import Button from '@/components/Button.vue'
import EventTitle from '@/components/EventTitle.vue'
import Input from '@/components/Input.vue'
import IconSend from '@/components/icons/IconSend.vue'
import type { EventDetail, MyResponse, Comment } from '@/service/events'
import { upsertEventResponse, deleteEvent, cancelEvent, uncancelEvent, listComments, createComment } from '@/service/events'
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

const comments = ref<Comment[]>([])
const newComment = ref('')
const isPostingComment = ref(false)
const COMMENT_MAX_LENGTH = 500

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
  loadComments()
})

const loadEvent = async () => {
  const response = await axios.get(`/events/${eventId}`)
  event.value = response.data
}

const loadComments = async () => {
  comments.value = await listComments(eventId)
}

const postComment = async () => {
  const body = newComment.value.trim()
  if (!body || isPostingComment.value) return
  isPostingComment.value = true
  try {
    const comment = await createComment(eventId, body)
    comments.value = [...comments.value, comment]
    newComment.value = ''
  } finally {
    isPostingComment.value = false
  }
}

// A comment's timestamp shown as a short, friendly absolute date and time.
const formatCommentTime = (isoTime: string) => {
  const date = new Date(isoTime)
  const datePart = date.toLocaleDateString(navigator.language, {
    month: 'short',
    day: 'numeric',
  })
  const timePart = date.toLocaleTimeString(navigator.language, {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,
  })
  return `${datePart}, ${timePart}`
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
    year: 'numeric',
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

          <!-- Comments -->
          <Divider class="eventDetailDivider" />
          <div class="eventDetailComments">
            <h2 class="eventDetailCommentsTitle">{{ i18n.global.t('events.detail.comments.title') }}</h2>

            <p v-if="comments.length === 0" class="eventDetailNoComments">
              {{ i18n.global.t('events.detail.comments.empty') }}
            </p>
            <div v-else class="eventDetailCommentList">
              <div v-for="comment in comments" :key="comment.id" class="eventDetailComment">
                <Avatar
                  :picture="comment.author.picture"
                  :given-name="comment.author.givenName"
                  :family-name="comment.author.familyName"
                  :nick-name="comment.author.nickName"
                  size="sm"
                />
                <div class="eventDetailCommentBody">
                  <div class="eventDetailCommentMeta">
                    <span class="eventDetailCommentAuthor">{{ comment.author.nickName || comment.author.givenName }}</span>
                    <span class="eventDetailCommentTime">{{ formatCommentTime(comment.createdAt) }}</span>
                  </div>
                  <p class="eventDetailCommentText">{{ comment.body }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Comment composer: pinned to the bottom of the screen on phones (the content
         above scrolls independently), a card below the event on wider screens. -->
    <div v-if="event" class="eventDetailComposer">
      <div class="eventDetailCommentForm">
        <Input
          id="eventDetailCommentInput"
          type="text"
          multi-line
          auto-grow
          :max-length="COMMENT_MAX_LENGTH"
          :placeholder="i18n.global.t('events.detail.comments.placeholder')"
          v-model="newComment"
        />
        <div class="eventDetailCommentSend">
          <Button
            icon-only
            :title="i18n.global.t('events.detail.comments.post')"
            :loading="isPostingComment"
            :disabled="newComment.trim().length === 0"
            @click="postComment"
          >
            <template #icon>
              <IconSend />
            </template>
          </Button>
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
  display: flex;
  flex-direction: column;
}

.eventDetailScroll {
  position: relative;
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
}

/* The composer is a flex sibling below the scroll area, so it stays pinned to the
   bottom of the screen while the event details and comments scroll above it — and it
   stays put even with no comments (unlike a sticky element, which needs content to
   push against). */
.eventDetailComposer {
  flex-shrink: 0;
  padding: var(--padding);
  background: var(--background-color);
  border-top: 1px solid var(--gray);
  box-sizing: border-box;
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

.eventDetailComments {
  display: flex;
  flex-direction: column;
  gap: var(--padding);
}

.eventDetailCommentsTitle {
  margin: 0;
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.eventDetailNoComments {
  color: var(--text-gray);
}

.eventDetailCommentList {
  display: flex;
  flex-direction: column;
  gap: var(--padding);
}

.eventDetailComment {
  display: flex;
  gap: var(--gap);
  align-items: flex-start;
}

.eventDetailCommentBody {
  flex: 1;
  min-width: 0;
}

.eventDetailCommentMeta {
  display: flex;
  align-items: baseline;
  gap: var(--gap);
}

.eventDetailCommentAuthor {
  font-weight: var(--font-weight-medium);
}

.eventDetailCommentTime {
  color: var(--text-gray);
  font-size: var(--font-size-small);
}

.eventDetailCommentText {
  margin-top: 2px;
  line-height: 1.4;
  overflow-wrap: anywhere;
  white-space: pre-wrap;
}

.eventDetailCommentForm {
  display: flex;
  align-items: flex-end;
  gap: var(--gap);
}

/* Let the input take the remaining width next to the fixed send button, and strip the
   floating-label's reserved top space so the bar sits low and compact. */
.eventDetailCommentForm :deep(.inputWrapper) {
  flex: 1;
  min-width: 0;
  padding-top: 0;
}

.eventDetailCommentForm :deep(.inputPlaceholder) {
  top: 10px;
}

/* On focus/with content the label would slide up into the stripped-away space, so keep
   it hidden and rely on the placeholder text alone — chat-bar style. */
.eventDetailCommentForm :deep(.input:is(:focus, :valid) ~ .inputPlaceholder) {
  display: none;
}

/* A round send button that stays a fixed size, so it never wobbles as you type. The
   wrapper is a real parent element so the :deep override reaches the Button's root. */
.eventDetailCommentSend {
  flex-shrink: 0;
}

.eventDetailCommentSend :deep(.button) {
  width: 44px;
  height: 44px;
  padding: 0;
  border-radius: 50%;
}

@media (min-width: 768px) {
  /* Desktop is a centered card the page scrolls around, so drop the mobile full-height
     flex layout and let the composer flow as its own card just below the event card. */
  .eventDetailPage {
    display: block;
    height: auto;
    padding: var(--padding);
    box-sizing: border-box;
  }

  .eventDetailScroll {
    position: relative;
    flex: none;
    min-height: 0;
    overflow-y: visible;
    max-width: var(--content-max-width);
    margin: 0 auto;
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
  }

  .eventDetailComposer {
    max-width: var(--content-max-width);
    margin: var(--padding) auto 0;
    border: none;
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

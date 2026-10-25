<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { isAxiosError } from 'axios'
import Modal from '@/components/Modal.vue'
import Button from '@/components/Button.vue'
import IconCircle from '@/components/IconCircle.vue'
import IconBell from '@/components/icons/IconBell.vue'
import RequestError from '@/components/RequestError.vue'
import { useAuthStore } from '@/stores/auth'
import { getPushState, turnOnPush } from '@/utils/push'

// Announces push notifications once per device. It only shows up while this device can turn
// them on (or on an iPhone that first needs the app on its home screen) — anyone who already
// has them on, blocked them, or uses a browser without support never sees it.
//
// The answer is remembered per device in localStorage, since notifications are per device
// too: "No thanks" hides it for good, "Later" brings it back after LATER_DAYS.
const STORAGE_KEY = 'notificationsAnnouncement'
const DISMISSED = 'dismissed'
const LATER_DAYS = 3

const auth = useAuthStore()

const step = ref<'intro' | 'homeScreen' | 'done' | null>(null)
const isTurningOn = ref(false)
// Browser-side failures; API errors are shown by RequestError.
const deviceError = ref(false)

const homeScreenSteps = [1, 2, 3, 4]

onMounted(async () => {
  if (!isDue()) return
  try {
    const state = await getPushState()
    if (state === 'off') step.value = 'intro'
    if (state === 'needsHomeScreen') step.value = 'homeScreen'
  } catch {
    // Not knowing the device's state is no reason to bother anyone with a modal.
  }
})

function isDue(): boolean {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved === DISMISSED) return false
    return saved === null || Date.now() >= Number(saved)
  } catch {
    return false
  }
}

function remember(value: string) {
  try {
    localStorage.setItem(STORAGE_KEY, value)
  } catch {
    // Without storage the modal may show up again, which is fine.
  }
}

async function turnOn() {
  isTurningOn.value = true
  deviceError.value = false
  try {
    const state = await turnOnPush(auth.user.email)
    if (state === 'on') {
      remember(DISMISSED)
      step.value = 'done'
    }
  } catch (e) {
    if (!isAxiosError(e)) {
      deviceError.value = true
    }
  } finally {
    isTurningOn.value = false
  }
}

function later() {
  remember(String(Date.now() + LATER_DAYS * 24 * 60 * 60 * 1000))
  step.value = null
}

function noThanks() {
  remember(DISMISSED)
  step.value = null
}
</script>

<template>
  <Modal v-if="step">
    <div class="notificationsAnnouncement">
      <IconCircle><IconBell /></IconCircle>

      <template v-if="step === 'done'">
        <h2 class="notificationsAnnouncementTitle">{{ $t('notificationsAnnouncement.done.title') }}</h2>
        <p class="notificationsAnnouncementText">{{ $t('notificationsAnnouncement.done.message') }}</p>
        <div class="notificationsAnnouncementActions">
          <Button :title="$t('notificationsAnnouncement.done.close')" @click="step = null" />
        </div>
      </template>

      <template v-else>
        <h2 class="notificationsAnnouncementTitle">{{ $t('notificationsAnnouncement.title') }}</h2>
        <p class="notificationsAnnouncementText">{{ $t('notificationsAnnouncement.message') }}</p>

        <template v-if="step === 'homeScreen'">
          <p class="notificationsAnnouncementText">{{ $t('notificationsAnnouncement.homeScreen.intro') }}</p>
          <ol class="notificationsAnnouncementSteps">
            <li v-for="n in homeScreenSteps" :key="n" class="notificationsAnnouncementStep">
              <span class="notificationsAnnouncementStepNumber">{{ n }}</span>
              <span class="notificationsAnnouncementStepText">{{ $t(`notificationsAnnouncement.homeScreen.step${n}`) }}</span>
            </li>
          </ol>
        </template>

        <p class="notificationsAnnouncementHint">{{ $t('notificationsAnnouncement.settingsHint') }}</p>
        <p v-if="deviceError" class="notificationsAnnouncementHint">{{ $t('notificationsAnnouncement.deviceError') }}</p>
        <RequestError class="notificationsAnnouncementError" />

        <div class="notificationsAnnouncementActions">
          <Button
            v-if="step === 'intro'"
            :title="$t('notificationsAnnouncement.turnOn')"
            :loading="isTurningOn"
            @click="turnOn"
          />
          <Button
            :title="$t('notificationsAnnouncement.noThanks')"
            theme="secondary"
            :disabled="isTurningOn"
            @click="noThanks"
          />
          <Button
            :title="$t('notificationsAnnouncement.later')"
            theme="tertiary"
            :disabled="isTurningOn"
            @click="later"
          />
        </div>
      </template>
    </div>
  </Modal>
</template>

<style scoped>
.notificationsAnnouncement {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap);
  text-align: center;
}

.notificationsAnnouncementTitle {
  margin: var(--gap) 0 0 0;
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.notificationsAnnouncementText {
  margin: 0;
  line-height: 1.5;
}

.notificationsAnnouncementHint {
  margin: 0;
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

.notificationsAnnouncementSteps {
  list-style: none;
  margin: var(--gap) 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--padding);
  text-align: left;
}

.notificationsAnnouncementStep {
  display: flex;
  align-items: flex-start;
  gap: var(--padding);
}

.notificationsAnnouncementStepNumber {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--gradient);
  color: var(--white);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
}

.notificationsAnnouncementStepText {
  padding-top: 4px;
  font-size: var(--font-size-small);
  line-height: 1.5;
}

.notificationsAnnouncementError {
  width: 100%;
  text-align: left;
}

.notificationsAnnouncementActions {
  width: 100%;
  margin-top: var(--padding);
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}
</style>

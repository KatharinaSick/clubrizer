<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useDocumentVisibility } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/Button.vue'
import RequestError from '@/components/RequestError.vue'
import kidsService from '@/service/kids'
import router from '@/router'
import { isAxiosError } from 'axios'
import { getPushState, turnOnPush, type PushState } from '@/utils/push'

const auth = useAuthStore()
const isChecking = ref(false)
const checked = ref(false)

// While waiting for approval, show the kids that were submitted so the parent can see
// their family is part of the request under review. Only fetched for pending accounts —
// a rejected account has nothing to look at (and the backend wouldn't return it).
const kidNames = ref<string[]>([])
const kidNamesLine = computed(() => kidNames.value.join(', '))

// Offer push notifications so the applicant hears about the decision without checking back.
// Only shown while it can actually be turned on (or explains the iPhone home-screen step);
// blocked or unsupported browsers just don't see the offer.
const pushState = ref<PushState | null>(null)
const isTurningOnPush = ref(false)
const pushDeviceError = ref(false)

// Check for a decision whenever the screen is opened or the app comes back to the front —
// e.g. after tapping the "you're approved" notification — so nobody has to tap "Check
// status" to find out.
async function refreshStatus() {
  await auth.checkStatus()
  if (auth.user.status === 'approved') {
    router.push('/')
  }
}

const visibility = useDocumentVisibility()
watch(visibility, (v) => {
  if (v === 'visible' && auth.user.status === 'pending') {
    refreshStatus()
  }
})

onMounted(async () => {
  if (auth.user.status !== 'pending') return
  refreshStatus()
  getPushState().then(s => { pushState.value = s }).catch(() => {})
  try {
    const kids = await kidsService.listKids()
    kidNames.value = kids.map((k) => k.givenName ?? '').filter(Boolean)
  } catch {
    // A missing kids line is not worth surfacing an error over on the waiting screen.
  }
})

async function checkStatus() {
  isChecking.value = true
  checked.value = false
  await refreshStatus()
  isChecking.value = false
  checked.value = auth.user.status === 'pending'
}

async function turnOnNotifications() {
  isTurningOnPush.value = true
  pushDeviceError.value = false
  try {
    pushState.value = await turnOnPush(auth.user.email)
  } catch (e) {
    // API errors are shown globally; browser-side failures are shown here.
    if (!isAxiosError(e)) {
      pushDeviceError.value = true
    }
  } finally {
    isTurningOnPush.value = false
  }
}

async function cancel() {
  await auth.logout()
  router.push('/signin')
}
</script>

<template>
  <div class="pendingApprovalContainer">
    <div class="pendingApprovalCenter">
      <img :alt="$t('team')" class="pendingApprovalLogo" src="@/assets/logo.svg" />
      <h1 class="pendingApprovalTitle">{{ $t('team') }}</h1>
    </div>

    <div class="pendingApprovalCenter">
      <template v-if="auth.user.status === 'rejected'">
        <h1>{{ $t('pendingApproval.rejected.title') }}</h1>
        <p>{{ $t('pendingApproval.rejected.message') }}</p>
      </template>
      <template v-else>
        <h1>{{ $t('pendingApproval.pending.title') }}</h1>
        <p>{{ $t('pendingApproval.pending.message') }}</p>
        <p v-if="kidNames.length > 0" class="pendingApprovalKids">
          {{ $t('pendingApproval.pending.yourKids', { names: kidNamesLine }) }}
        </p>
      </template>
    </div>

    <div v-if="auth.user.status !== 'rejected'" class="pendingApprovalCenter">
      <Button
        :title="$t('pendingApproval.pending.checkStatus')"
        :loading="isChecking"
        theme="secondary"
        @click="checkStatus"
      />
      <p v-if="checked" class="pendingApprovalStillPending">{{ $t('pendingApproval.pending.stillPending') }}</p>
    </div>

    <div
      v-if="auth.user.status === 'pending' && (pushState === 'off' || pushState === 'on' || pushState === 'needsHomeScreen')"
      class="pendingApprovalCenter pendingApprovalNotifications"
    >
      <template v-if="pushState === 'off'">
        <p class="pendingApprovalNotificationsText">{{ $t('pendingApproval.notifications.question') }}</p>
        <Button
          :title="$t('pendingApproval.notifications.turnOn')"
          theme="tertiary"
          accent
          :loading="isTurningOnPush"
          @click="turnOnNotifications"
        />
      </template>
      <p v-else-if="pushState === 'on'" class="pendingApprovalNotificationsText">
        {{ $t('pendingApproval.notifications.on') }}
      </p>
      <template v-else>
        <p class="pendingApprovalNotificationsText">{{ $t('pendingApproval.notifications.question') }}</p>
        <p class="pendingApprovalStillPending">{{ $t('pendingApproval.notifications.homeScreen') }}</p>
      </template>
      <p v-if="pushDeviceError" class="pendingApprovalStillPending">{{ $t('pendingApproval.notifications.deviceError') }}</p>
      <RequestError class="pendingApprovalNotificationsError" />
    </div>

    <div class="pendingApprovalCenter pendingApprovalCancel">
      <Button :title="$t('pendingApproval.cancel')" theme="red" @click="cancel" />
    </div>
  </div>
</template>

<style scoped>
.pendingApprovalContainer {
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  padding-top: 64px;

  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 48px;
}

.pendingApprovalCancel {
  margin-top: auto;
}

.pendingApprovalCenter {
  text-align: center;
}

.pendingApprovalLogo {
  height: 100px;
}

.pendingApprovalTitle {
  margin: 8px 0 0 0;
  text-transform: uppercase;
  letter-spacing: .1rem;
  background-image: var(--gradient);
  color: transparent;
  background-clip: text;
  -webkit-background-clip: text;
}

.pendingApprovalKids {
  margin: 4px 0 0 0;
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

.pendingApprovalStillPending {
  margin: 12px 0 0 0;
  font-size: var(--font-size-small);
  color: var(--text-gray);
  text-align: center;
}

.pendingApprovalNotifications {
  max-width: 400px;
}

.pendingApprovalNotificationsText {
  margin: 0 0 var(--gap) 0;
}

.pendingApprovalNotificationsError {
  margin-top: var(--gap);
  text-align: left;
}

@media (min-width: 768px) {
  .pendingApprovalContainer {
    justify-content: center;
    padding-top: 0;
    max-width: 400px;
    margin: 0 auto;
  }

  .pendingApprovalCancel {
    margin-top: -24px;
  }
}
</style>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useDocumentVisibility } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/Button.vue'
import RequestError from '@/components/RequestError.vue'
import IconBell from '@/components/icons/IconBell.vue'
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
async function refreshStatus(silent = true) {
  await auth.checkStatus(silent)
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
  await refreshStatus(false)
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
    <img :alt="$t('team')" class="pendingApprovalLogo" src="@/assets/logo.svg" />

    <div class="pendingApprovalHero">
      <template v-if="auth.user.status === 'rejected'">
        <h1 class="pendingApprovalTitle">{{ $t('pendingApproval.rejected.title') }}</h1>
        <p class="pendingApprovalMessage">{{ $t('pendingApproval.rejected.message') }}</p>
      </template>
      <template v-else>
        <h1 class="pendingApprovalTitle">{{ $t('pendingApproval.pending.title') }}</h1>
        <p class="pendingApprovalMessage">{{ $t('pendingApproval.pending.message') }}</p>
        <p v-if="kidNames.length > 0" class="pendingApprovalHint">
          {{ $t('pendingApproval.pending.yourKids', { names: kidNamesLine }) }}
        </p>
      </template>
    </div>

    <template v-if="auth.user.status !== 'rejected'">
      <div
        v-if="auth.user.status === 'pending' && (pushState === 'off' || pushState === 'on' || pushState === 'needsHomeScreen')"
        class="pendingApprovalNotifications"
      >
        <IconBell class="pendingApprovalNotificationsIcon" />
        <div class="pendingApprovalNotificationsBody">
          <p class="pendingApprovalNotificationsText">
            {{ pushState === 'on' ? $t('pendingApproval.notifications.on') : $t('pendingApproval.notifications.question') }}
          </p>
          <p v-if="pushState === 'needsHomeScreen'" class="pendingApprovalHint">
            {{ $t('pendingApproval.notifications.homeScreen') }}
          </p>
          <p v-if="pushDeviceError" class="pendingApprovalHint">{{ $t('pendingApproval.notifications.deviceError') }}</p>
        </div>
        <Button
          v-if="pushState === 'off'"
          :title="$t('pendingApproval.notifications.turnOn')"
          theme="secondary"
          inline
          :loading="isTurningOnPush"
          @click="turnOnNotifications"
        />
      </div>

      <div class="pendingApprovalCheck">
        <Button
          :title="$t('pendingApproval.pending.checkStatus')"
          :loading="isChecking"
          theme="secondary"
          @click="checkStatus"
        />
        <p v-if="checked" class="pendingApprovalHint">{{ $t('pendingApproval.pending.stillPending') }}</p>
      </div>
    </template>

    <RequestError class="pendingApprovalError" />

    <div class="pendingApprovalSignOut">
      <Button :title="$t('pendingApproval.signOut')" theme="tertiary" @click="cancel" />
    </div>
  </div>
</template>

<style scoped>
.pendingApprovalContainer {
  width: 100%;
  min-height: 100%;
  max-width: 400px;
  margin: 0 auto;
  box-sizing: border-box;
  padding-top: calc(var(--padding) * 6);

  display: flex;
  flex-direction: column;
  align-items: center;
  gap: calc(var(--padding) * 3);
}

.pendingApprovalLogo {
  height: 80px;
  margin-bottom: calc(var(--padding) * 2);
}

.pendingApprovalHero {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--gap);
  text-align: center;
  margin-top: var(--padding);
}

.pendingApprovalTitle {
  margin: var(--gap) 0 0 0;
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.pendingApprovalMessage {
  margin: 0;
}

.pendingApprovalHint {
  margin: 0;
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

.pendingApprovalNotifications {
  width: 100%;
  box-sizing: border-box;
  padding: var(--padding);
  background: var(--light-gray);
  border-radius: var(--border-radius);
  display: flex;
  align-items: center;
  gap: var(--padding);
}

.pendingApprovalNotificationsIcon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  color: var(--blue);
}

.pendingApprovalNotificationsBody {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: calc(var(--gap) / 2);
}

.pendingApprovalNotificationsText {
  margin: 0;
  font-size: var(--font-size-small);
}

.pendingApprovalCheck {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  text-align: center;
}

.pendingApprovalError {
  width: 100%;
}

.pendingApprovalSignOut {
  margin-top: auto;
  padding-bottom: var(--padding);
}

@media (min-width: 768px) {
  .pendingApprovalContainer {
    justify-content: center;
    padding-top: 0;
  }

  .pendingApprovalLogo {
    margin-bottom: 0;
  }

  .pendingApprovalSignOut {
    margin-top: 0;
  }
}
</style>

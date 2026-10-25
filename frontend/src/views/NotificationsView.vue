<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { isAxiosError } from 'axios'
import Header from '@/components/Header.vue'
import Button from '@/components/Button.vue'
import Alert from '@/components/Alert.vue'
import RequestError from '@/components/RequestError.vue'
import { useAuthStore } from '@/stores/auth'
import notificationsService from '@/service/notifications'
import {
  getPushState,
  turnOffPushOnAllDevices,
  turnOffPushOnThisDevice,
  turnOnPush,
  type PushState,
} from '@/utils/push'

const router = useRouter()
const auth = useAuthStore()

// null until the device has been checked, so no state flashes up before it's known.
const state = ref<PushState | null>(null)
const busyAction = ref<'turnOn' | 'test' | 'offThisDevice' | 'offAllDevices' | null>(null)
const testSent = ref(false)
// Browser-side failures (e.g. the push service refusing to subscribe). API errors are
// shown globally through RequestError instead.
const deviceError = ref(false)

onMounted(() => run(null, async () => {
  state.value = await getPushState()
}))

// run wraps every action: it shows the button spinner, resets old feedback and sorts
// browser-side errors (shown locally) from API errors (shown by RequestError).
async function run(action: typeof busyAction.value, fn: () => Promise<void>) {
  busyAction.value = action
  testSent.value = false
  deviceError.value = false
  try {
    await fn()
  } catch (e) {
    if (!isAxiosError(e)) {
      deviceError.value = true
    }
  } finally {
    busyAction.value = null
  }
}

const turnOn = () => run('turnOn', async () => {
  state.value = await turnOnPush(auth.user.email)
})

const sendTest = () => run('test', async () => {
  await notificationsService.sendTest()
  testSent.value = true
})

const turnOffThisDevice = () => run('offThisDevice', async () => {
  await turnOffPushOnThisDevice()
  state.value = 'off'
})

const turnOffAllDevices = () => run('offAllDevices', async () => {
  await turnOffPushOnAllDevices()
  state.value = 'off'
})
</script>

<template>
  <div class="notificationsView">
    <div class="notificationsCard">
      <Header
        :title="$t('notifications.header')"
        left-action="back"
        :back-fn="() => router.push('/profile')"
      />

      <template v-if="state === 'needsHomeScreen'">
        <p class="notificationsText">{{ $t('notifications.homeScreen.intro') }}</p>
        <ol class="notificationsSteps">
          <li>{{ $t('notifications.homeScreen.step1') }}</li>
          <li>{{ $t('notifications.homeScreen.step2') }}</li>
          <li>{{ $t('notifications.homeScreen.step3') }}</li>
        </ol>
      </template>

      <Alert
        v-else-if="state === 'unsupported'"
        variant="info"
        :message="$t('notifications.unsupported')"
      />

      <Alert
        v-else-if="state === 'blocked'"
        variant="warning"
        :title="$t('notifications.blocked.title')"
        :message="$t('notifications.blocked.message')"
      />

      <template v-else-if="state === 'off'">
        <p class="notificationsText">{{ $t('notifications.off.intro') }}</p>
        <div class="notificationsActions">
          <Button
            :title="$t('notifications.off.turnOn')"
            :loading="busyAction === 'turnOn'"
            @click="turnOn"
          />
        </div>
      </template>

      <template v-else-if="state === 'on'">
        <Alert variant="info" :message="$t('notifications.on.intro')" />
        <div class="notificationsActions">
          <Button
            :title="$t('notifications.on.sendTest')"
            theme="secondary"
            :loading="busyAction === 'test'"
            @click="sendTest"
          />
          <template v-if="testSent">
            <p class="notificationsHint">{{ $t('notifications.on.testSent') }}</p>
            <p class="notificationsHint">{{ $t('notifications.on.testMissing') }}</p>
          </template>
          <Button
            :title="$t('notifications.on.turnOffThisDevice')"
            theme="tertiary"
            :loading="busyAction === 'offThisDevice'"
            @click="turnOffThisDevice"
          />
          <Button
            :title="$t('notifications.on.turnOffAllDevices')"
            theme="tertiary"
            danger
            :loading="busyAction === 'offAllDevices'"
            @click="turnOffAllDevices"
          />
        </div>
      </template>

      <Alert
        v-if="deviceError"
        class="notificationsError"
        :message="$t('notifications.deviceError')"
      />
      <RequestError class="notificationsError" />
    </div>
  </div>
</template>

<style scoped>
.notificationsText {
  margin: 0 0 var(--padding) 0;
}

.notificationsSteps {
  margin: 0;
  padding-left: calc(var(--padding) * 2);
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.notificationsActions {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  margin-top: calc(var(--padding) * 2);
}

.notificationsHint {
  margin: 0;
  text-align: center;
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

.notificationsError {
  margin-top: calc(var(--padding) * 2);
}

@media (min-width: 768px) {
  .notificationsView {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .notificationsCard {
    width: 100%;
    max-width: var(--content-max-width);
    background: var(--background-color);
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
    padding: var(--padding);
    box-sizing: border-box;
  }
}
</style>

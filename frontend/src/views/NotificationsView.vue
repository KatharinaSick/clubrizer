<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { isAxiosError } from 'axios'
import Header from '@/components/Header.vue'
import Button from '@/components/Button.vue'
import Alert from '@/components/Alert.vue'
import Divider from '@/components/Divider.vue'
import Toggle from '@/components/Toggle.vue'
import IconCircle from '@/components/IconCircle.vue'
import NumberedSteps from '@/components/NumberedSteps.vue'
import IconBell from '@/components/icons/IconBell.vue'
import RequestError from '@/components/RequestError.vue'
import i18n from '@/plugins/i18n'
import { useAuthStore } from '@/stores/auth'
import notificationsService, { type NotificationSetting } from '@/service/notifications'
import {
  getPushState,
  turnOffPushOnThisDevice,
  turnOnPush,
  type PushState,
} from '@/utils/push'

const router = useRouter()
const auth = useAuthStore()

// null until the device has been checked, so no state flashes up before it's known.
const state = ref<PushState | null>(null)
const busyAction = ref<'turnOn' | 'test' | 'offThisDevice' | null>(null)
const testSent = ref(false)
// The test reached none of the user's devices (e.g. the registration expired).
const testUnreachable = ref(false)
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
  testUnreachable.value = false
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

// Settings apply to all of the user's devices, but are only shown while this device has
// notifications on — otherwise the switches would have no visible effect here. Pending
// accounts only get their approval, which can't be switched off, so they see none.
const settings = ref<NotificationSetting[]>([])
const savingCategory = ref<string | null>(null)

watch(state, async (newState) => {
  if (newState === 'on' && auth.user.status === 'approved' && settings.value.length === 0) {
    settings.value = await notificationsService.getSettings()
  }
})

async function updateSetting(setting: NotificationSetting, enabled: boolean) {
  setting.enabled = enabled
  savingCategory.value = setting.category
  try {
    await notificationsService.updateSetting(setting.category, enabled)
  } catch {
    // Put the switch back; RequestError shows what went wrong.
    setting.enabled = !enabled
  } finally {
    savingCategory.value = null
  }
}

const homeScreenSteps = computed(() => [1, 2, 3, 4].map((n) => i18n.global.t(`notifications.homeScreen.step${n}`)))

const turnOn = () => run('turnOn', async () => {
  state.value = await turnOnPush(auth.user.email)
})

const sendTest = () => run('test', async () => {
  const delivered = await notificationsService.sendTest()
  testSent.value = delivered > 0
  testUnreachable.value = delivered === 0
})

const turnOffThisDevice = () => run('offThisDevice', async () => {
  await turnOffPushOnThisDevice()
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

      <div v-if="state === 'needsHomeScreen'" class="notificationsIntro">
        <IconCircle><IconBell /></IconCircle>
        <h2 class="notificationsIntroTitle">{{ $t('notifications.homeScreen.title') }}</h2>
        <p class="notificationsIntroText">{{ $t('notifications.homeScreen.intro') }}</p>
        <NumberedSteps class="notificationsSteps" :steps="homeScreenSteps" />
      </div>

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

      <div v-else-if="state === 'off'" class="notificationsIntro">
        <IconCircle><IconBell /></IconCircle>
        <h2 class="notificationsIntroTitle">{{ $t('notifications.off.title') }}</h2>
        <p class="notificationsIntroText">{{ $t('notifications.off.intro') }}</p>
        <Button
          class="notificationsIntroButton"
          :title="$t('notifications.off.turnOn')"
          :loading="busyAction === 'turnOn'"
          @click="turnOn"
        />
      </div>

      <template v-else-if="state === 'on'">
        <section class="notificationsSection">
          <div class="notificationsStatus">
            <span class="notificationsStatusDot" aria-hidden="true" />
            <span class="notificationsStatusText">{{ $t('notifications.on.status') }}</span>
            <Button
              :title="$t('notifications.on.sendTest')"
              theme="secondary"
              inline
              :loading="busyAction === 'test'"
              @click="sendTest"
            />
          </div>
          <template v-if="testSent">
            <p class="notificationsHint">{{ $t('notifications.on.testSent') }}</p>
            <p class="notificationsHint">{{ $t('notifications.on.testMissing') }}</p>
          </template>
          <p v-if="testUnreachable" class="notificationsHint">{{ $t('notifications.on.testUnreachable') }}</p>
        </section>

        <section v-if="settings.length > 0" class="notificationsSection">
          <h2 class="notificationsSectionTitle">{{ $t('notifications.settings.title') }}</h2>
          <div class="notificationsSettings">
            <template v-for="(setting, index) in settings" :key="setting.category">
              <Divider v-if="index > 0" />
              <Toggle
                :id="`notification-setting-${setting.category}`"
                :label="$t(`notifications.settings.${setting.category}.label`)"
                :description="$t(`notifications.settings.${setting.category}.description`)"
                :disabled="savingCategory === setting.category"
                :model-value="setting.enabled"
                @update:model-value="updateSetting(setting, $event)"
              />
            </template>
          </div>
        </section>

        <section class="notificationsSection">
          <div class="notificationsTurnOff">
            <Button
              :title="$t('notifications.on.turnOffThisDevice')"
              theme="secondary"
              :loading="busyAction === 'offThisDevice'"
              @click="turnOffThisDevice"
            />
          </div>
        </section>
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
.notificationsIntro {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding-top: calc(var(--padding) * 2);
}

.notificationsIntroTitle {
  margin: calc(var(--padding) * 1.5) 0 var(--gap) 0;
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.notificationsIntroText {
  margin: 0;
  max-width: 320px;
}

.notificationsIntroButton {
  margin-top: calc(var(--padding) * 2);
}

.notificationsSteps {
  margin-top: calc(var(--padding) * 2);
}

.notificationsSection + .notificationsSection {
  margin-top: calc(var(--padding) * 2);
}

.notificationsSectionTitle {
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-bold);
  margin: 0 0 var(--padding);
}

.notificationsStatus {
  padding: var(--padding);
  background: var(--light-green);
  border-radius: var(--border-radius);
  display: flex;
  align-items: center;
  gap: var(--gap);
}

.notificationsStatusDot {
  flex-shrink: 0;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--green);
}

.notificationsStatusText {
  flex: 1;
  font-weight: var(--font-weight-medium);
}

.notificationsSettings {
  display: flex;
  flex-direction: column;
  gap: var(--padding);
}

.notificationsTurnOff {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.notificationsHint {
  margin: var(--gap) 0 0 0;
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

  .notificationsIntroButton {
    width: auto;
    padding-inline: calc(var(--padding) * 3);
  }

  .notificationsTurnOff {
    flex-direction: row;
    justify-content: flex-end;
  }

  .notificationsTurnOff :deep(.button) {
    width: auto;
  }
}
</style>

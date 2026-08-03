<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRequestStore } from '@/stores/request'
import Button from '@/components/Button.vue'
import OnboardingKids from '@/components/OnboardingKids.vue'
import RequestError from '@/components/RequestError.vue'
import kidsService from '@/service/kids'
import i18n from '@/plugins/i18n'
import router from '@/router'

type KidType = 'meAndKids' | 'onlyKids'
type BusyAction = KidType | 'justMe' | 'submit' | null

const auth = useAuthStore()
const requestStore = useRequestStore()

const step = ref<'type' | 'kids'>('type')
const initialNames = ref<string[]>([])
const validationError = ref('')
const busyAction = ref<BusyAction>(null)
const busy = computed(() => busyAction.value !== null)
const kidsRef = ref<InstanceType<typeof OnboardingKids> | null>(null)

onMounted(async () => {
  try {
    const kids = await kidsService.listKids()
    // Seed the name fields from any kids already saved (a rare onboarding interrupted
    // after a partial save), so nothing the parent typed is lost.
    initialNames.value = kids.map((k) => k.givenName ?? '').filter(Boolean)
    // Resume the kids step if the account already committed to a kid path (guardian,
    // or a participant that already added kids); otherwise start at the choice.
    if (!auth.user.selfParticipates || kids.length > 0) {
      step.value = 'kids'
    }
  } catch {
    // errors surface via RequestError
  }
})

async function ensureAccountType(selfParticipates: boolean) {
  if (auth.user.selfParticipates !== selfParticipates) {
    await auth.setAccountType(selfParticipates)
  }
}

async function finish() {
  await auth.finishOnboarding()
  router.push('/pending-approval')
}

// "Just me" needs nothing more — submit straight to approval.
async function chooseJustMe() {
  if (busy.value) return
  busyAction.value = 'justMe'
  requestStore.clearError()
  try {
    await ensureAccountType(true)
    await finish()
  } catch {
    // errors surface via RequestError
  } finally {
    busyAction.value = null
  }
}

// The kid options advance to the kids step instead of submitting.
async function chooseKidType(type: KidType) {
  if (busy.value) return
  busyAction.value = type
  requestStore.clearError()
  try {
    await ensureAccountType(type === 'meAndKids')
    step.value = 'kids'
  } catch {
    // errors surface via RequestError
  } finally {
    busyAction.value = null
  }
}

function onKidsChange(count: number) {
  if (count > 0) validationError.value = ''
}

async function submitKids() {
  if (busy.value) return
  validationError.value = ''
  requestStore.clearError()
  busyAction.value = 'submit'
  try {
    const names = kidsRef.value?.names() ?? []
    // Both kid paths ("me & my kids" and "only my kids") require at least one kid — this
    // step is only reached by choosing a kid option, so an empty list is always a mistake.
    // (An account that wants no kids picks "just me", which skips this step entirely.)
    if (names.length === 0) {
      validationError.value = i18n.global.t('onboarding.needKid')
      return
    }
    // Save the whole list in one request (replace semantics), then submit for approval.
    await kidsService.replaceKids(names)
    await finish()
  } catch {
    // errors surface via RequestError
  } finally {
    busyAction.value = null
  }
}

function back() {
  validationError.value = ''
  requestStore.clearError()
  step.value = 'type'
}

async function cancel() {
  await auth.logout()
  router.push('/signin')
}
</script>

<template>
  <div class="onboardingContainer">
    <div class="onboardingCenter">
      <img :alt="$t('team')" class="onboardingLogo" src="@/assets/logo.svg" />
      <h1 class="onboardingTitle">{{ $t('team') }}</h1>
    </div>

    <div class="onboardingBody">
      <!-- Step 1: pick who the account is for. Each option is an action. -->
      <template v-if="step === 'type'">
        <h2 class="onboardingSectionTitle">{{ $t('onboarding.title') }}</h2>
        <p class="onboardingSectionHint">{{ $t('onboarding.hint') }}</p>
        <RequestError class="onboardingRequestError" />
        <div class="onboardingTypeOptions">
          <Button
            :title="$t('onboarding.justMe')"
            theme="secondary"
            :loading="busyAction === 'justMe'"
            :disabled="busy"
            @click="chooseJustMe"
          />
          <Button
            :title="$t('onboarding.meAndKids')"
            theme="secondary"
            :loading="busyAction === 'meAndKids'"
            :disabled="busy"
            @click="chooseKidType('meAndKids')"
          />
          <Button
            :title="$t('onboarding.onlyKids')"
            theme="secondary"
            :loading="busyAction === 'onlyKids'"
            :disabled="busy"
            @click="chooseKidType('onlyKids')"
          />
        </div>
      </template>

      <!-- Step 2: enter the kids' first names. -->
      <template v-else>
        <h2 class="onboardingSectionTitle">{{ $t('onboarding.kidsTitle') }}</h2>
        <p class="onboardingSectionHint">{{ $t('onboarding.kidsHint') }}</p>
        <OnboardingKids ref="kidsRef" :initial-names="initialNames" @change="onKidsChange" />
        <p v-if="validationError" class="onboardingValidationError">{{ validationError }}</p>
        <RequestError class="onboardingRequestError" />
        <div class="onboardingKidsActions">
          <Button
            :title="$t('onboarding.continue')"
            :loading="busyAction === 'submit'"
            @click="submitKids"
          />
          <Button
            :title="$t('onboarding.back')"
            theme="tertiary"
            :disabled="busy"
            @click="back"
          />
        </div>
      </template>
    </div>

    <div class="onboardingCenter onboardingCancel">
      <Button :title="$t('onboarding.cancel')" theme="tertiary" :disabled="busy" @click="cancel" />
    </div>
  </div>
</template>

<style scoped>
.onboardingContainer {
  width: 100%;
  min-height: 100%;
  box-sizing: border-box;
  padding: 64px var(--padding) var(--padding);

  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 32px;
}

.onboardingCenter {
  text-align: center;
}

.onboardingCancel {
  margin-top: auto;
}

.onboardingLogo {
  height: 100px;
}

.onboardingTitle {
  margin: 8px 0 0 0;
  text-transform: uppercase;
  letter-spacing: 0.1rem;
  background-image: var(--gradient);
  color: transparent;
  background-clip: text;
  -webkit-background-clip: text;
}

.onboardingBody {
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.onboardingSectionTitle {
  margin: 0;
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-bold);
  text-align: center;
}

.onboardingSectionHint {
  margin: 0 0 var(--gap);
  font-size: var(--font-size-small);
  color: var(--text-gray);
  text-align: center;
}

.onboardingTypeOptions {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.onboardingKidsActions {
  margin-top: var(--padding);
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.onboardingValidationError {
  color: var(--red);
  font-size: var(--font-size-small);
}

.onboardingRequestError:empty {
  display: none;
}

@media (min-width: 768px) {
  .onboardingContainer {
    justify-content: flex-start;
    padding-top: var(--padding);
    max-width: 400px;
    margin: 0 auto;
  }
}
</style>

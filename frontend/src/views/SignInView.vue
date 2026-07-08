<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRequestStore } from '@/stores/request'
import { useRoute } from 'vue-router'
import router from '@/router'
import Input from '@/components/Input.vue'
import Button from '@/components/Button.vue'
import Header from '@/components/Header.vue'
import RequestError from '@/components/RequestError.vue'
import i18n from '@/plugins/i18n'

const auth = useAuthStore()
const requestStore = useRequestStore()
const route = useRoute()

type Step = 'email' | 'code'
const step = ref<Step>('email')
const email = ref('')
const code = ref('')
const emailError = ref('')
const codeError = ref('')
const isLoading = ref(false)
const isResending = ref(false)
const codeSent = ref(false)

async function submitEmail() {
  if (isLoading.value) return
  emailError.value = ''
  if (!email.value) {
    emailError.value = i18n.global.t('signIn.emailRequired')
    return
  }
  isLoading.value = true
  try {
    await auth.requestOTP(email.value)
    step.value = 'code'
  } catch {
    // error shown globally via RequestError
  } finally {
    isLoading.value = false
  }
}

async function submitCode() {
  if (isLoading.value) return
  codeError.value = ''
  if (code.value.length !== 6) {
    codeError.value = i18n.global.t('signIn.codeInvalid')
    return
  }
  isLoading.value = true
  try {
    await auth.verifyOTP(email.value, code.value)
    router.replace({ path: (route.query.redirect as string) ?? '/' })
  } catch {
    // error shown globally via RequestError
  } finally {
    isLoading.value = false
  }
}

async function resendCode() {
  isResending.value = true
  codeSent.value = false
  codeError.value = ''
  code.value = ''
  requestStore.clearError()
  try {
    await auth.requestOTP(email.value)
    codeSent.value = true
  } catch {
    // error shown globally via RequestError
  } finally {
    isResending.value = false
  }
}

function backToEmail() {
  step.value = 'email'
  code.value = ''
  codeError.value = ''
  codeSent.value = false
  requestStore.clearError()
}
</script>

<template>
  <div class="signInContainer">
    <div v-if="step === 'code'" class="signInCodeHeader">
      <Header left-action="back" :back-fn="backToEmail" />
    </div>

    <div class="signInCard">

      <template v-if="step === 'email'">
        <div class="signInHeader">
          <h1 class="signInWelcome">{{ $t('signIn.welcomeTo') }}<br />{{ $t('team') }}!</h1>
        </div>
        <form class="signInForm" @submit.prevent="submitEmail">
          <p class="signInSubtitle">{{ $t('signIn.getStarted') }}</p>
          <Input
            id="email"
            type="email"
            :placeholder="$t('signIn.emailLabel')"
            v-model="email"
            :error="emailError"
            required
            theme="ghost"
          />
          <Button :title="$t('signIn.sendCode')" :loading="isLoading" theme="ghost" />
        </form>
        <p class="signInNote">{{ $t('signIn.emailNote') }}</p>
        <RouterLink to="/getting-started" class="signInHowItWorks">{{ $t('signIn.howItWorks') }}</RouterLink>
        <RequestError />
      </template>

      <template v-else>
        <div class="signInHeader">
          <h1 class="signInWelcome">{{ $t('signIn.codeTitle') }}</h1>
        </div>
        <form class="signInForm" @submit.prevent="submitCode">
          <p class="signInSubtitle">{{ $t('signIn.codeNote', { email }) }}</p>
          <Input
            id="code"
            type="text"
            inputMode="numeric"
            :maxLength="6"
            :placeholder="$t('signIn.codeLabel')"
            v-model="code"
            :error="codeError"
            required
            theme="ghost"
          />
          <Button :title="$t('signIn.verify')" :loading="isLoading" theme="ghost" />
        </form>
        <RequestError />
        <div class="signInResend">
          <p v-if="codeSent" class="signInCodeSent">{{ $t('signIn.codeSent') }}</p>
          <Button :title="$t('signIn.resendCode')" :loading="isResending" :disabled="isLoading" theme="tertiary" @click="resendCode" />
        </div>
      </template>

    </div>
  </div>
</template>

<style scoped>
.signInContainer {
  position: fixed;
  inset: 0;
  z-index: 10;
  background: linear-gradient(45deg, var(--gradient-start-light), var(--gradient-end-light));

  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px var(--padding);
  box-sizing: border-box;
  overflow-y: auto;
}

.signInCard {
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

.signInHeader {
  text-align: center;
}

.signInWelcome {
  margin: 0;
  font-size: 28px;
  font-weight: var(--font-weight-bold);
  color: var(--white);
  line-height: 1.2;
}

.signInSubtitle {
  margin: 0 0 var(--gap) 0;
  color: var(--white);
  font-size: var(--font-size-medium);
}

.signInForm {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.signInNote {
  margin: 0;
  text-align: center;
  color: var(--white);
  opacity: 0.7;
  font-size: var(--font-size-small);
}

.signInHowItWorks {
  color: var(--white);
  opacity: 0.9;
  font-size: var(--font-size-small);
  text-decoration: underline;
}

.signInHowItWorks:active {
  opacity: 0.6;
}

.signInCodeHeader {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  padding: 4px var(--padding);
}

.signInCodeHeader :deep(.headerIcon) {
  color: var(--white);
}

.signInResend :deep(.button.tertiary) {
  color: var(--white);
  opacity: 0.8;
}

.signInResend {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.signInCodeSent {
  margin: 0;
  font-size: var(--font-size-small);
  color: var(--white);
  opacity: 0.8;
  text-align: center;
}
</style>

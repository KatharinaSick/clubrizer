<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/Button.vue'
import router from '@/router'

const auth = useAuthStore()
const isChecking = ref(false)
const checked = ref(false)

async function checkStatus() {
  isChecking.value = true
  checked.value = false
  await auth.refreshTokens()
  isChecking.value = false
  if (auth.user.status === 'approved') {
    router.push('/')
  } else {
    checked.value = true
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
      <img
        alt="Clubrizer Logo"
        class="pendingApprovalLogo"
        src="@/assets/logo.svg"
      />
      <h1 class="pendingApprovalTitle">Clubrizer</h1>
      <h3 class="pendingApprovalSlogan">Team Up!</h3>
    </div>

    <div class="pendingApprovalCenter">
      <template v-if="auth.user.status === 'rejected'">
        <h1>{{ $t('pendingApproval.rejected.title') }}</h1>
        <p>{{ $t('pendingApproval.rejected.message') }}</p>
      </template>
      <template v-else>
        <h1>{{ $t('pendingApproval.pending.title') }}</h1>
        <p>{{ $t('pendingApproval.pending.message') }}</p>
      </template>
    </div>

    <div v-if="auth.user.status !== 'rejected'" class="pendingApprovalCenter">
      <Button :title="$t('pendingApproval.pending.checkStatus')" :loading="isChecking" theme="secondary" @click="checkStatus" />
      <p v-if="checked" class="pendingApprovalStillPending">{{ $t('pendingApproval.pending.stillPending') }}</p>
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
  margin-top: 64px;

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
}

.pendingApprovalSlogan {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: .1rem;
  background-image: var(--gradient);
  color: transparent;
  background-clip: text;
}

.pendingApprovalStillPending {
  margin: 12px 0 0 0;
  font-size: var(--font-size-small);
  color: var(--text-gray);
  text-align: center;
}
</style>

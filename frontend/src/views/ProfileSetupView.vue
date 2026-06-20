<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRequestStore } from '@/stores/request'
import router from '@/router'
import i18n from '@/plugins/i18n'
import Header from '@/components/Header.vue'
import Input from '@/components/Input.vue'
import Button from '@/components/Button.vue'
import RequestError from '@/components/RequestError.vue'

const auth = useAuthStore()
const requestStore = useRequestStore()

const firstName = ref('')
const lastName = ref('')
const nickName = ref('')
const picture = ref<File | null>(null)
const isLoading = ref(false)
const firstNameError = ref('')
const lastNameError = ref('')

function onPictureChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  picture.value = file ?? null
}

async function submit() {
  firstNameError.value = ''
  lastNameError.value = ''
  requestStore.clearError()

  if (!firstName.value) {
    firstNameError.value = i18n.global.t('profileSetup.firstNameRequired')
    return
  }
  if (!lastName.value) {
    lastNameError.value = i18n.global.t('profileSetup.lastNameRequired')
    return
  }

  isLoading.value = true
  try {
    if (picture.value) {
      await auth.updateProfilePicture(picture.value)
    }
    await auth.updateProfile(firstName.value, lastName.value, nickName.value || undefined)
    router.replace('/events')
  } catch {
    // errors shown globally via RequestError
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="profileSetupView">
    <div class="profileSetupCard">
      <Header :title="$t('profileSetup.header')" show-divider />

      <div class="profileSetupNameRow">
        <Input
          id="firstName"
          type="text"
          :placeholder="$t('profileSetup.firstName')"
          v-model="firstName"
          :error="firstNameError"
          required
        />
        <Input
          id="lastName"
          type="text"
          :placeholder="$t('profileSetup.lastName')"
          v-model="lastName"
          :error="lastNameError"
          required
        />
      </div>
      <Input
        id="nickName"
        type="text"
        :placeholder="$t('profileSetup.nickName')"
        v-model="nickName"
      />

      <div class="profileSetupFileWrapper">
        <label for="picture" class="profileSetupFileLabel">{{ $t('profileSetup.picture') }}</label>
        <input
          id="picture"
          type="file"
          accept="image/*"
          class="profileSetupFileInput"
          @change="onPictureChange"
        />
      </div>

      <RequestError class="profileSetupError" />

      <Button
        :title="$t('profileSetup.save')"
        :loading="isLoading"
        @click="submit"
        class="profileSetupButton"
      />
    </div>
  </div>
</template>

<style scoped>
.profileSetupNameRow {
  display: flex;
  flex-direction: column;
}

.profileSetupFileWrapper {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.profileSetupFileLabel {
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

.profileSetupFileInput {
  font-size: var(--font-size-small);
  font-family: inherit;
  color: var(--text-light);
}

.profileSetupError {
  margin-top: 24px;
}

.profileSetupButton {
  margin-top: 24px;
}

@media (min-width: 768px) {
  .profileSetupView {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .profileSetupCard {
    width: 100%;
    max-width: var(--content-max-width);
    background: var(--background-color);
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
    padding: var(--padding);
    box-sizing: border-box;
  }

  .profileSetupNameRow {
    flex-direction: row;
    gap: var(--gap);
  }

  .profileSetupButton {
    width: auto;
    padding-inline: calc(var(--padding) * 3);
    display: block;
    margin-left: auto;
  }
}
</style>

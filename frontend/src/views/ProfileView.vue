<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRequestStore } from '@/stores/request'
import { useRouter } from 'vue-router'
import ProfileInfo from '@/components/ProfileInfo.vue'
import Header from '@/components/Header.vue'
import Button from '@/components/Button.vue'
import Input from '@/components/Input.vue'
import RequestError from '@/components/RequestError.vue'
import i18n from '@/plugins/i18n'

const auth = useAuthStore()
const requestStore = useRequestStore()
const router = useRouter()

const isEditing = ref(false)
const isLoading = ref(false)
const firstName = ref('')
const lastName = ref('')
const nickName = ref('')
const picture = ref<File | null>(null)
const firstNameError = ref('')
const lastNameError = ref('')

function startEditing() {
  firstName.value = auth.user.givenName ?? ''
  lastName.value = auth.user.familyName ?? ''
  nickName.value = auth.user.nickName ?? ''
  picture.value = null
  firstNameError.value = ''
  lastNameError.value = ''
  requestStore.clearError()
  isEditing.value = true
}

function cancelEditing() {
  isEditing.value = false
  requestStore.clearError()
}

function onPictureChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  picture.value = file ?? null
}

async function saveProfile() {
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
    isEditing.value = false
  } catch {
    // errors shown globally via RequestError
  } finally {
    isLoading.value = false
  }
}

const logout = async () => {
  await auth.logout()
  router.push('/signin')
}
</script>

<template>
  <div>
    <Header :title="i18n.global.t('profile.header')" />

    <template v-if="!isEditing">
      <ProfileInfo :user="auth.user" />
      <div class="profileActions">
        <Button :title="$t('profile.edit')" theme="secondary" @click="startEditing" />
        <Button :title="$t('profile.logout')" theme="red" @click="logout" />
      </div>
    </template>

    <template v-else>
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
      <Input
        id="nickName"
        type="text"
        :placeholder="$t('profileSetup.nickName')"
        v-model="nickName"
      />
      <div class="profileEditFileWrapper">
        <label for="picture" class="profileEditFileLabel">{{ $t('profileSetup.picture') }}</label>
        <input
          id="picture"
          type="file"
          accept="image/*"
          class="profileEditFileInput"
          @change="onPictureChange"
        />
      </div>
      <RequestError class="profileEditError" />
      <div class="profileActions">
        <Button :title="$t('profileSetup.save')" :loading="isLoading" @click="saveProfile" />
        <Button :title="$t('profile.cancel')" theme="secondary" :disabled="isLoading" @click="cancelEditing" />
      </div>
    </template>
  </div>
</template>

<style scoped>
.profileActions {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  margin-top: 24px;
}

.profileEditFileWrapper {
  margin-top: 24px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.profileEditFileLabel {
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

.profileEditFileInput {
  font-size: var(--font-size-small);
  font-family: inherit;
  color: var(--text-light);
}

.profileEditError {
  margin-top: 24px;
}
</style>

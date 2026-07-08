<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRequestStore } from '@/stores/request'
import { useRouter } from 'vue-router'
import ProfileInfo from '@/components/ProfileInfo.vue'
import Header from '@/components/Header.vue'
import Button from '@/components/Button.vue'
import Input from '@/components/Input.vue'
import RequestError from '@/components/RequestError.vue'
import i18n from '@/plugins/i18n'
import usersService, { type Role } from '@/service/users'

const auth = useAuthStore()
const requestStore = useRequestStore()
const router = useRouter()

const roles = ref<Role[]>([])

onMounted(() => {
  usersService.getMyRoles().then(r => { roles.value = r }).catch(() => {})
})

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
  <div class="profileView">
    <div class="profileCard">
      <Header :title="i18n.global.t('profile.header')" />

      <template v-if="!isEditing">
        <ProfileInfo :user="auth.user" :roles="roles" />
        <div class="profileActions">
          <Button :title="$t('profile.edit')" theme="secondary" @click="startEditing" />
          <Button :title="$t('profile.logout')" theme="red" @click="logout" />
        </div>
        <div class="profileHowItWorks">
          <RouterLink to="/getting-started" class="profileHowItWorksLink">{{ $t('profile.howItWorks') }}</RouterLink>
        </div>
      </template>

      <template v-else>
        <div class="profileNameRow">
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
  </div>
</template>

<style scoped>
.profileNameRow {
  display: flex;
  flex-direction: column;
}

.profileActions {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  margin-top: 24px;
}

.profileHowItWorks {
  margin-top: var(--padding);
  text-align: center;
}

.profileHowItWorksLink {
  color: var(--text-gray);
  font-size: var(--font-size-small);
  text-decoration: underline;
}

.profileHowItWorksLink:active {
  opacity: 0.6;
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

@media (min-width: 768px) {
  .profileView {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .profileCard {
    width: 100%;
    max-width: var(--content-max-width);
    background: var(--background-color);
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
    padding: var(--padding);
    box-sizing: border-box;
  }

  .profileNameRow {
    flex-direction: row;
    gap: var(--gap);
  }

  .profileActions {
    flex-direction: row;
    justify-content: flex-end;
  }

  .profileActions :deep(.button) {
    width: auto;
    padding-inline: calc(var(--padding) * 3);
  }
}
</style>

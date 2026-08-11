<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useApprovalsStore } from '@/stores/approvals'
import { useRequestStore } from '@/stores/request'
import { useRouter } from 'vue-router'
import ProfileInfo from '@/components/ProfileInfo.vue'
import Header from '@/components/Header.vue'
import Button from '@/components/Button.vue'
import Input from '@/components/Input.vue'
import MenuButton, { type MenuItem } from '@/components/MenuButton.vue'
import IconMore from '@/components/icons/IconMore.vue'
import RequestError from '@/components/RequestError.vue'
import i18n from '@/plugins/i18n'
import usersService, { type Role } from '@/service/users'

const auth = useAuthStore()
const approvals = useApprovalsStore()
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

// Profile actions live in the gear menu. "Manage Members" and its red dot only appear for
// accounts that can manage members; everyone else sees the personal actions only.
const menuItems = computed<MenuItem[]>(() => {
  const items: MenuItem[] = [
    { label: i18n.global.t('profile.edit'), onClick: startEditing },
    { label: i18n.global.t('profile.myKids'), onClick: () => router.push('/my-kids') },
  ]
  if (auth.canManageMembers) {
    items.push({
      label: i18n.global.t('profile.manageMembers'),
      badge: approvals.hasPending,
      onClick: () => router.push('/manage-members'),
    })
  }
  items.push({ label: i18n.global.t('profile.logout'), danger: true, onClick: logout })
  items.push({ label: i18n.global.t('profile.howItWorks'), subtle: true, divider: true, onClick: () => router.push('/getting-started') })
  items.push({ label: i18n.global.t('profile.changelog'), subtle: true, onClick: () => router.push('/changelog') })
  return items
})
</script>

<template>
  <div class="profileView">
    <div class="profileCard">
      <Header
        :title="i18n.global.t('profile.header')"
        :left-action="isEditing ? 'back' : undefined"
        :back-fn="isEditing ? cancelEditing : undefined"
      >
        <template #right>
          <MenuButton
            v-if="!isEditing"
            variant="bare"
            placement="top-aligned"
            :items="menuItems"
            :badge="auth.canManageMembers && approvals.hasPending"
            :aria-label="i18n.global.t('profile.menuLabel')"
          >
            <template #icon><IconMore class="profileMenuIcon" /></template>
          </MenuButton>
        </template>
      </Header>

      <template v-if="!isEditing">
        <ProfileInfo :user="auth.user" :roles="roles" />
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

.profileMenuIcon {
  width: 22px;
  height: 22px;
  color: var(--text-color);
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

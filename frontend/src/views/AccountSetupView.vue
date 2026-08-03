<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRequestStore } from '@/stores/request'
import router from '@/router'
import i18n from '@/plugins/i18n'
import Header from '@/components/Header.vue'
import Avatar from '@/components/Avatar.vue'
import Divider from '@/components/Divider.vue'
import Input from '@/components/Input.vue'
import Button from '@/components/Button.vue'
import RequestError from '@/components/RequestError.vue'
import kidsService from '@/service/kids'

// One-time post-approval setup. Every approved account fills in at least a first and last
// name (guardians included — for club records — but without nickname/photo). Accounts with
// kids also give each kid a last name (photo optional). There's no "done" flag: the router
// gates on the own name being present, and submit saves the own name LAST, so setup can't
// be half-finished.

const auth = useAuthStore()
const requestStore = useRequestStore()

const participates = computed(() => auth.user.selfParticipates)

const firstName = ref(auth.user.givenName ?? '')
const lastName = ref(auth.user.familyName ?? '')
const nickName = ref(auth.user.nickName ?? '')
const ownPicture = ref<File | null>(null)
const ownPreview = ref<string | null>(null)
const firstNameError = ref('')
const lastNameError = ref('')

type KidDraft = {
  id: string
  givenName: string
  lastName: string
  lastNameError: string
  picture: File | null
  preview: string | null
  existingPicture: string | null
}
const kidDrafts = ref<KidDraft[]>([])

const isLoading = ref(false)

onMounted(async () => {
  try {
    const kids = await kidsService.listKids()
    kidDrafts.value = kids
      .map((k) => ({
        id: k.id,
        givenName: k.givenName ?? '',
        lastName: k.familyName ?? '',
        lastNameError: '',
        picture: null,
        preview: null,
        existingPicture: k.picture ?? null,
      }))
  } catch {
    // errors surface via RequestError
  }
})

function onOwnPictureChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0] ?? null
  ownPicture.value = file
  ownPreview.value = file ? URL.createObjectURL(file) : null
}

function onKidPictureChange(draft: KidDraft, event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0] ?? null
  draft.picture = file
  draft.preview = file ? URL.createObjectURL(file) : null
}

function validate(): boolean {
  let ok = true
  firstNameError.value = ''
  lastNameError.value = ''
  if (!firstName.value) {
    firstNameError.value = i18n.global.t('accountSetup.firstNameRequired')
    ok = false
  }
  if (!lastName.value) {
    lastNameError.value = i18n.global.t('accountSetup.lastNameRequired')
    ok = false
  }
  for (const draft of kidDrafts.value) {
    draft.lastNameError = ''
    if (!draft.lastName) {
      draft.lastNameError = i18n.global.t('accountSetup.kidLastNameRequired')
      ok = false
    }
  }
  return ok
}

async function submit() {
  requestStore.clearError()
  if (!validate()) return

  isLoading.value = true
  try {
    // Kids first, own profile last — saving the own name is what lifts the setup gate, so
    // it must only happen once every required field is in. A mid-way exit re-triggers setup
    // rather than leaving kids half-done.
    for (const draft of kidDrafts.value) {
      await kidsService.updateKid(draft.id, draft.givenName, draft.lastName)
      if (draft.picture) {
        await kidsService.updateKidPicture(draft.id, draft.picture)
      }
    }

    if (participates.value && ownPicture.value) {
      await auth.updateProfilePicture(ownPicture.value)
    }
    await auth.updateProfile(firstName.value, lastName.value, nickName.value || undefined)

    router.replace('/events')
  } catch {
    // errors surface via RequestError
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="accountSetupView">
    <div class="accountSetupCard">
      <Header :title="$t('accountSetup.header')" show-divider />
      <p class="accountSetupIntro">{{ $t('accountSetup.intro') }}</p>

      <!-- Own data: one titled block, mirroring each kid's block below. -->
      <h2 class="accountSetupSectionTitle">{{ $t('accountSetup.yourData') }}</h2>
      <div class="accountSetupPerson">
        <label v-if="participates" class="accountSetupAvatarButton" :aria-label="$t('accountSetup.changePhoto')">
          <Avatar
            :picture="ownPreview ?? auth.user.picture"
            :given-name="firstName || null"
            :family-name="lastName || null"
            :nick-name="nickName || null"
            size="lg"
          />
          <span class="accountSetupAvatarOverlay">📷</span>
          <input type="file" accept="image/*" class="accountSetupHiddenFile" @change="onOwnPictureChange" />
        </label>
        <div class="accountSetupFields">
          <div class="accountSetupNameRow">
            <Input
              id="firstName"
              type="text"
              :placeholder="$t('accountSetup.firstName')"
              v-model="firstName"
              :error="firstNameError"
              required
            />
            <Input
              id="lastName"
              type="text"
              :placeholder="$t('accountSetup.lastName')"
              v-model="lastName"
              :error="lastNameError"
              required
            />
          </div>
          <Input
            v-if="participates"
            id="nickName"
            type="text"
            :placeholder="$t('accountSetup.nickName')"
            v-model="nickName"
          />
        </div>
      </div>

      <!-- One titled block per kid. -->
      <template v-for="draft in kidDrafts" :key="draft.id">
        <Divider class="accountSetupDivider" />
        <h2 class="accountSetupSectionTitle">{{ draft.givenName }}</h2>
        <div class="accountSetupPerson accountSetupPersonKid">
          <label class="accountSetupAvatarButton" :aria-label="$t('accountSetup.kidPhoto')">
            <Avatar
              :picture="draft.preview ?? draft.existingPicture"
              :given-name="draft.givenName || null"
              :family-name="draft.lastName || null"
              size="lg"
            />
            <span class="accountSetupAvatarOverlay">📷</span>
            <input type="file" accept="image/*" class="accountSetupHiddenFile" @change="onKidPictureChange(draft, $event)" />
          </label>
          <div class="accountSetupFields">
            <Input
              :id="`kidLast-${draft.id}`"
              type="text"
              :placeholder="$t('accountSetup.kidLastName')"
              v-model="draft.lastName"
              :error="draft.lastNameError"
              required
            />
          </div>
        </div>
      </template>

      <RequestError class="accountSetupError" />

      <Button
        :title="$t('accountSetup.save')"
        :loading="isLoading"
        @click="submit"
        class="accountSetupButton"
      />
    </div>
  </div>
</template>

<style scoped>
.accountSetupIntro {
  color: var(--text-gray);
  font-size: var(--font-size-small);
  line-height: 1.4;
  margin-bottom: var(--padding);
}

.accountSetupSectionTitle {
  margin: var(--padding) 0 0;
  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-bold);
}

.accountSetupDivider {
  margin-top: var(--padding);
}

.accountSetupPerson {
  display: flex;
  gap: var(--padding);
  align-items: flex-start;
  margin-top: var(--gap);
}

/* A kid block is a single field next to the avatar. Nudge the avatar down so it lines up
   with the input's visible box (the input reserves ~24px on top for its floating label),
   rather than centring against the whole wrapper — which left the field looking low. */
.accountSetupPersonKid .accountSetupAvatarButton {
  margin-top: var(--padding);
}

.accountSetupFields {
  flex: 1;
  min-width: 0;
}

.accountSetupNameRow {
  display: flex;
  flex-direction: column;
}

.accountSetupAvatarButton {
  position: relative;
  flex-shrink: 0;
  margin-top: calc(var(--padding) * 2);
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  border-radius: 50%;
  display: flex;
}

/* Prominent bottom-right camera badge so "tap to add a photo" reads clearly even before a
   photo is set (the Avatar shows initials underneath). */
.accountSetupAvatarOverlay {
  position: absolute;
  bottom: 2px;
  right: 2px;
  font-size: var(--font-size-large);
  line-height: 1;
}

.accountSetupHiddenFile {
  display: none;
}

.accountSetupError {
  margin-top: var(--padding);
}

.accountSetupButton {
  margin-top: var(--padding);
}

@media (min-width: 768px) {
  .accountSetupView {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .accountSetupCard {
    width: 100%;
    max-width: var(--content-max-width);
    background: var(--background-color);
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
    padding: var(--padding);
    box-sizing: border-box;
  }

  .accountSetupNameRow {
    flex-direction: row;
    gap: var(--gap);
  }

  .accountSetupButton {
    width: auto;
    padding-inline: calc(var(--padding) * 3);
    display: block;
    margin-left: auto;
  }
}
</style>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import kidsService, { type Kid } from '@/service/kids'
import i18n from '@/plugins/i18n'
import Avatar from '@/components/Avatar.vue'
import Button from '@/components/Button.vue'
import Input from '@/components/Input.vue'
import Modal from '@/components/Modal.vue'
import RequestError from '@/components/RequestError.vue'
import IconPencil from '@/components/icons/IconPencil.vue'
import IconTrash from '@/components/icons/IconTrash.vue'

// Post-approval kid management: a calm read-only list (photo, name, status) with edit and
// remove actions, plus one "add" action. Add and edit share a single modal that takes all
// three fields at once — photo, first name, last name — so the section stays consistent
// with the rest of the app. The lightweight name-only onboarding entry lives in
// OnboardingKids.vue.

const kids = ref<Kid[]>([])

// One modal drives both add and edit; mode decides the title and which save path runs.
const modalMode = ref<'add' | 'edit' | null>(null)
const editingKid = ref<Kid | null>(null)
const formFirstName = ref('')
const formLastName = ref('')
const firstNameError = ref('')
const lastNameError = ref('')
const formPicture = ref<File | null>(null)
const formPreview = ref<string | null>(null)
const existingPicture = ref<string | null>(null)
const isSaving = ref(false)

const removeTarget = ref<Kid | null>(null)
const isRemoving = ref(false)

onMounted(load)

async function load() {
  try {
    kids.value = await kidsService.listKids()
  } catch {
    // errors shown globally via RequestError
  }
}

function resetForm() {
  formFirstName.value = ''
  formLastName.value = ''
  firstNameError.value = ''
  lastNameError.value = ''
  if (formPreview.value) URL.revokeObjectURL(formPreview.value)
  formPicture.value = null
  formPreview.value = null
  existingPicture.value = null
}

function openAdd() {
  resetForm()
  editingKid.value = null
  modalMode.value = 'add'
}

function openEdit(kid: Kid) {
  resetForm()
  editingKid.value = kid
  formFirstName.value = kid.givenName ?? ''
  formLastName.value = kid.familyName ?? ''
  existingPicture.value = kid.picture ?? null
  modalMode.value = 'edit'
}

function closeModal() {
  if (isSaving.value) return
  modalMode.value = null
  editingKid.value = null
}

function onPictureChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0] ?? null
  if (formPreview.value) URL.revokeObjectURL(formPreview.value)
  formPicture.value = file
  formPreview.value = file ? URL.createObjectURL(file) : null
}

function validate(): boolean {
  firstNameError.value = ''
  lastNameError.value = ''
  let ok = true
  if (!formFirstName.value.trim()) {
    firstNameError.value = i18n.global.t('kids.firstNameRequired')
    ok = false
  }
  if (!formLastName.value.trim()) {
    lastNameError.value = i18n.global.t('kids.lastNameRequired')
    ok = false
  }
  return ok
}

async function save() {
  if (isSaving.value || !validate()) return
  isSaving.value = true
  try {
    if (modalMode.value === 'add') {
      let kid = await kidsService.addKid(formFirstName.value.trim(), formLastName.value.trim())
      // The picture endpoint needs an existing kid, so upload it right after creation.
      if (formPicture.value) kid = await kidsService.updateKidPicture(kid.id, formPicture.value)
      kids.value.push(kid)
    } else if (editingKid.value) {
      const id = editingKid.value.id
      let kid = await kidsService.updateKid(id, formFirstName.value.trim(), formLastName.value.trim())
      if (formPicture.value) kid = await kidsService.updateKidPicture(id, formPicture.value)
      const index = kids.value.findIndex((k) => k.id === id)
      if (index !== -1) kids.value[index] = kid
    }
    modalMode.value = null
    editingKid.value = null
  } catch {
    // errors shown globally via RequestError
  } finally {
    isSaving.value = false
  }
}

async function confirmRemove() {
  if (!removeTarget.value || isRemoving.value) return
  isRemoving.value = true
  try {
    await kidsService.removeKid(removeTarget.value.id)
    kids.value = kids.value.filter((k) => k.id !== removeTarget.value!.id)
    removeTarget.value = null
  } catch {
    removeTarget.value = null
  } finally {
    isRemoving.value = false
  }
}
</script>

<template>
  <div class="kidsManager">
    <RequestError class="kidsManagerError" />

    <ul v-if="kids.length > 0" class="kidsManagerList">
      <li v-for="kid in kids" :key="kid.id" class="kidsManagerKid">
        <Avatar :picture="kid.picture" :given-name="kid.givenName" :family-name="kid.familyName" size="lg" />

        <div class="kidsManagerBody">
          <span class="kidsManagerName">{{ kid.givenName }} {{ kid.familyName }}</span>
          <span v-if="kid.status === 'pending'" class="kidsManagerBadge kidsManagerBadgePending">
            {{ $t('kids.pendingBadge') }}
          </span>
        </div>

        <div class="kidsManagerActions">
          <Button inline icon-only accent theme="tertiary" :title="$t('kids.edit')" @click="openEdit(kid)">
            <template #icon><IconPencil /></template>
          </Button>
          <Button inline icon-only danger theme="tertiary" :title="$t('kids.remove')" @click="removeTarget = kid">
            <template #icon><IconTrash /></template>
          </Button>
        </div>
      </li>
    </ul>
    <p v-else class="kidsManagerEmpty">{{ $t('kids.empty') }}</p>

    <Button :title="$t('kids.add')" theme="secondary" @click="openAdd" />

    <!-- Add & edit share this modal: photo, first name, last name in one place. -->
    <Modal v-if="modalMode">
      <div class="kidsManagerForm">
        <h2 class="kidsManagerFormTitle">
          {{ modalMode === 'add' ? $t('kids.add') : $t('kids.editTitle') }}
        </h2>

        <label class="kidsManagerFormAvatar" :aria-label="$t('kids.changePhoto')">
          <Avatar
            :picture="formPreview ?? existingPicture"
            :given-name="formFirstName || null"
            :family-name="formLastName || null"
            size="lg"
          />
          <span class="kidsManagerFormCam">📷</span>
          <input type="file" accept="image/*" class="kidsManagerFormFile" @change="onPictureChange" />
        </label>

        <div class="kidsManagerFormFields">
          <Input
            id="kidFirstName"
            type="text"
            :placeholder="$t('kids.firstName')"
            v-model="formFirstName"
            :error="firstNameError"
            required
            @keyup.enter="save"
          />
          <Input
            id="kidLastName"
            type="text"
            :placeholder="$t('kids.lastName')"
            v-model="formLastName"
            :error="lastNameError"
            required
            @keyup.enter="save"
          />
        </div>

        <div class="kidsManagerFormActions">
          <Button :title="$t('kids.save')" :loading="isSaving" @click="save" />
          <Button :title="$t('kids.cancel')" theme="secondary" :disabled="isSaving" @click="closeModal" />
        </div>
      </div>
    </Modal>

    <Modal v-if="removeTarget">
      <div class="kidsManagerRemoveConfirm">
        <h2 class="kidsManagerRemoveTitle">{{ $t('kids.removeConfirm.title') }}</h2>
        <p class="kidsManagerRemoveMessage">
          {{ $t('kids.removeConfirm.message', { name: removeTarget.givenName }) }}
        </p>
        <div class="kidsManagerRemoveActions">
          <Button :title="$t('kids.removeConfirm.confirm')" theme="red" :loading="isRemoving" @click="confirmRemove" />
          <Button :title="$t('kids.removeConfirm.cancel')" theme="secondary" :disabled="isRemoving" @click="removeTarget = null" />
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
.kidsManager {
  display: flex;
  flex-direction: column;
  gap: var(--padding);
}

.kidsManagerList {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
}

.kidsManagerKid {
  display: flex;
  align-items: center;
  gap: var(--padding);
  padding: var(--padding) 0;
}

.kidsManagerKid + .kidsManagerKid {
  border-top: 1px solid var(--gray);
}

.kidsManagerBody {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.kidsManagerName {
  font-weight: var(--font-weight-medium);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kidsManagerBadge {
  align-self: flex-start;
  font-size: var(--font-size-small);
  padding: 2px var(--gap);
  border-radius: var(--border-radius);
}

.kidsManagerBadgePending {
  background: var(--light-gray);
  color: var(--text-gray);
}

.kidsManagerActions {
  display: flex;
  gap: var(--gap);
  align-items: center;
  flex-shrink: 0;
}

.kidsManagerEmpty {
  color: var(--text-gray);
  font-size: var(--font-size-small);
}

/* Shared add/edit form inside the modal. */
.kidsManagerForm {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--padding);
}

.kidsManagerFormTitle {
  margin: 0;
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
  align-self: flex-start;
}

.kidsManagerFormAvatar {
  position: relative;
  flex-shrink: 0;
  display: flex;
  padding: 0;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 50%;
}

/* Prominent bottom-right camera badge so "tap to add a photo" reads clearly even before a
   photo is set (the Avatar shows initials underneath). */
.kidsManagerFormCam {
  position: absolute;
  bottom: 2px;
  right: 2px;
  font-size: var(--font-size-large);
  line-height: 1;
}

.kidsManagerFormFile {
  display: none;
}

.kidsManagerFormFields {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.kidsManagerFormActions {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  margin-top: var(--gap);
}

.kidsManagerRemoveConfirm {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.kidsManagerRemoveTitle {
  margin: 0;
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.kidsManagerRemoveMessage {
  margin: 0;
  line-height: 1.5;
}

.kidsManagerRemoveActions {
  margin-top: var(--padding);
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

@media (min-width: 768px) {
  .kidsManagerFormFields {
    flex-direction: row;
  }
}
</style>

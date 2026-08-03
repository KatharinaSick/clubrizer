<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import i18n from '@/plugins/i18n'
import Button from '@/components/Button.vue'
import Input from '@/components/Input.vue'

// Lightweight given-name entry for onboarding: an auto-growing column of first-name
// fields with no per-row "add" button. A trailing empty field always sits at the bottom;
// typing into it spawns the next one, so a parent just types names down the list. Nothing
// is persisted here — the parent screen collects names() on submit and saves them in one
// request. The richer post-approval manager lives in KidsManager.vue.
const props = defineProps<{ initialNames?: string[] }>()
const emit = defineEmits<{ change: [count: number] }>()

type Row = { key: number; name: string }
let nextKey = 0
const rows = ref<Row[]>([])

for (const name of props.initialNames ?? []) {
  rows.value.push({ key: nextKey++, name })
}
appendEmptyRow()

function appendEmptyRow() {
  rows.value.push({ key: nextKey++, name: '' })
}

function removeRow(key: number) {
  rows.value = rows.value.filter((r) => r.key !== key)
  if (rows.value.length === 0) appendEmptyRow()
}

function names(): string[] {
  return rows.value.map((r) => r.name.trim()).filter(Boolean)
}

const hasAnyName = computed(() => rows.value.some((r) => r.name.trim() !== ''))

// The first field reads as the primary "First name" input. Once a name is entered, the
// auto-added trailing field is labelled as an optional "add another" so it never looks
// like a required blank the parent must fill before continuing.
function placeholderFor(index: number): string {
  const isTrailingAddField = index === rows.value.length - 1 && hasAnyName.value
  return isTrailingAddField
    ? i18n.global.t('onboarding.addAnotherKid')
    : i18n.global.t('onboarding.kidNamePlaceholder')
}

// Keep exactly one trailing empty field, and report the current non-empty count so the
// parent screen can clear its validation hint as soon as a name is entered.
watch(
  rows,
  (value) => {
    const last = value[value.length - 1]
    if (!last || last.name.trim() !== '') appendEmptyRow()
    emit('change', names().length)
  },
  { deep: true },
)

defineExpose({ names })
</script>

<template>
  <div class="onboardingKids">
    <div v-for="(row, index) in rows" :key="row.key" class="onboardingKidsRow">
      <Input
        :id="`onboardingKid-${row.key}`"
        type="text"
        :placeholder="placeholderFor(index)"
        v-model="row.name"
      />
      <Button
        v-if="row.name.trim() !== '' || index !== rows.length - 1"
        inline
        icon-only
        theme="tertiary"
        class="onboardingKidsRemove"
        :title="$t('onboarding.removeKid')"
        @click="removeRow(row.key)"
      >
        <template #icon>✕</template>
      </Button>
    </div>
  </div>
</template>

<style scoped>
.onboardingKids {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.onboardingKidsRow {
  display: flex;
  align-items: center;
  gap: var(--gap);
}

.onboardingKidsRow :deep(.inputWrapper) {
  flex: 1;
}

.onboardingKidsRemove {
  flex-shrink: 0;
  margin-top: var(--padding);
}
</style>

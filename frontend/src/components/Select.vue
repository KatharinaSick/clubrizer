<script setup lang="ts">
import i18n from '../plugins/i18n'

export interface Option {
  name: string
  id: string
}

defineProps<{
  id: string
  placeholder: string
  options: Option[]
  loading?: boolean
}>()

const value = defineModel()
</script>

<template>
  <select
    :id="id"
    class="select"
    v-model="value"
  >
    <option :value="null" disabled selected>{{ placeholder }}</option>
    <option v-if="loading" disabled>{{ i18n.global.t('select.loading') }}</option>
    <option v-else v-for="option in options" :value="option">{{ option.name }}</option>
  </select>
</template>

<style scoped>
.select {
  width: 100%;

  margin-top: 24px;
  padding: var(--padding-input);
  font: inherit;

  border-radius: var(--border-radius);
  border: none;
  box-shadow: var(--box-shadow);

  &::-ms-expand {
    display: none;
  }

  &:focus {
    outline: none;
  }
}
</style>

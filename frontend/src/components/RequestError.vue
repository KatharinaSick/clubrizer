<script setup lang="ts">
import { useRequestStore } from '@/stores/request'
import Alert from '@/components/Alert.vue'
import i18n from '@/plugins/i18n'

const requestStore = useRequestStore()
</script>

<template>
  <Alert
    v-if="requestStore.error"
    :title="i18n.global.t('request.errorTitle')"
    :message="requestStore.error"
    variant="error"
  />
</template>
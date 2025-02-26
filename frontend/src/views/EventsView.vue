<script setup lang="ts">
import Header from '@/components/Header.vue'
import axios from '@/plugins/axios'
import i18n from '@/plugins/i18n'

const sendRequest = () => {
  axios.get('/events').then(response => console.log(response))
}
</script>

<template>
  <div>
    <Header :title="i18n.global.t('events.header')" />

  <button @click="sendRequest">Send request</button>
  </div>
</template>

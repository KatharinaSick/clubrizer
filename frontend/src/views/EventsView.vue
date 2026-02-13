<script setup lang="ts">
import Header from '@/components/Header.vue'
import axios from '@/plugins/axios'
import i18n from '@/plugins/i18n'
import FloatingActionButton from '@/components/FloatingActionButton.vue'
import { onMounted, ref } from 'vue'
import type { EventProps } from '@/components/Event.vue'
import Event from '@/components/Event.vue'

const eventsLoading = ref(true)
const eventsError = ref(null)
const events = ref<EventProps[]>([])

onMounted(() => {
  // TODO show loading indicator
  eventsLoading.value = true
  eventsError.value = null

  axios
    .get('/events')
    .then(response => {
      eventsLoading.value = false
      events.value = response.data
      console.log(events.value)
    })
    .catch(error => {
      // TODO handle
      console.log(error)
    })
})
</script>

<template>
  <div>
    <Header :title="i18n.global.t('events.header')" />

    <Event v-for="event in events" :event="event" />

    <RouterLink to="/events/new"><FloatingActionButton /></RouterLink>
  </div>
</template>

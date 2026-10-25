<script setup lang="ts">
import Header from '@/components/Header.vue'
import EventForm from '@/components/EventForm.vue'
import type { EventFormValues } from '@/components/EventForm.vue'
import RequestError from '@/components/RequestError.vue'
import i18n from '@/plugins/i18n'
import router from '@/router'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getEvent, updateEvent } from '@/service/events'

const route = useRoute()
const eventId = route.params.id as string

const initialValues = ref<EventFormValues | null>(null)
const saveLoading = ref(false)

onMounted(async () => {
  const event = await getEvent(eventId)
  initialValues.value = {
    title: event.title,
    description: event.description,
    startTime: new Date(event.startTime),
    location: event.location,
  }
})

const save = (values: EventFormValues) => {
  saveLoading.value = true

  updateEvent(eventId, values)
    .then(() => {
      saveLoading.value = false
      router.back()
    })
    .catch(() => {
      saveLoading.value = false
    })
}
</script>

<template>
  <div class="editEventView">
    <div class="editEventCard">
      <Header
        :title="i18n.global.t('events.edit.header')"
        show-divider
        left-action="back"
      />

      <EventForm
        v-if="initialValues"
        :submit-title="i18n.global.t('events.edit.save')"
        :loading="saveLoading"
        :initial-values="initialValues"
        @submit="save"
      />
      <RequestError v-else />
    </div>
  </div>
</template>

<style scoped>
@media (min-width: 768px) {
  .editEventView {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .editEventCard {
    width: 100%;
    max-width: var(--content-max-width);
    background: var(--background-color);
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
    padding: var(--padding);
    box-sizing: border-box;
  }
}
</style>

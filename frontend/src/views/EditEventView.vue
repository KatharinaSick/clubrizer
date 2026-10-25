<script setup lang="ts">
import Alert from '@/components/Alert.vue'
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
const notEditable = ref(false)
const saveLoading = ref(false)

onMounted(async () => {
  let event
  try {
    event = await getEvent(eventId)
  } catch {
    // The error is shown by <RequestError />.
    return
  }
  // The edit page can also be opened from a link or the browser history, so check here
  // too that the event may still be edited. The backend enforces this as well.
  if (!event.canEdit || event.cancelledAt || new Date(event.startTime) <= new Date()) {
    notEditable.value = true
    return
  }
  initialValues.value = {
    title: event.title,
    description: event.description,
    startTime: new Date(event.startTime),
    location: event.location,
  }
})

// Go back to the event after saving. If the edit page was opened directly, there is no
// previous page in the app to go back to, so open the event instead.
const leave = () => {
  if (window.history.state?.back) {
    router.back()
  } else {
    router.replace({ name: 'event-detail', params: { id: eventId } })
  }
}

const save = async (values: EventFormValues) => {
  saveLoading.value = true
  try {
    await updateEvent(eventId, values)
    leave()
  } catch {
    // The error is shown by <RequestError /> inside the form.
  } finally {
    saveLoading.value = false
  }
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
      <Alert
        v-else-if="notEditable"
        :message="i18n.global.t('events.edit.notEditable')"
        variant="info"
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

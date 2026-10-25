<script setup lang="ts">
import Header from '@/components/Header.vue'
import EventForm from '@/components/EventForm.vue'
import type { EventFormValues } from '@/components/EventForm.vue'
import i18n from '@/plugins/i18n'
import router from '@/router'
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { createEvent as createEventApi } from '@/service/events'

const route = useRoute()
const categoryId = route.params.categoryId as string
const categoryName = route.query.categoryName as string

const createLoading = ref(false)

const createEvent = (values: EventFormValues) => {
  createLoading.value = true;

  createEventApi({ ...values, categoryId })
    .then(event => {
      createLoading.value = false;
      router.replace(`/events/${event.id}`)
    })
    .catch(() => {
      createLoading.value = false;
    })
}
</script>

<template>
  <div class="newEventView">
    <div class="newEventCard">
      <Header
        :title="i18n.global.t('events.new.header')"
        show-divider
        left-action="back"
      />

      <EventForm
        :submit-title="i18n.global.t('events.new.create')"
        :loading="createLoading"
        :initial-title="categoryName"
        @submit="createEvent"
      />
    </div>
  </div>
</template>

<style scoped>
@media (min-width: 768px) {
  .newEventView {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .newEventCard {
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

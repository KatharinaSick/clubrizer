<script setup lang="ts">
import Header from '@/components/Header.vue'
import i18n from '@/plugins/i18n'
import router from '@/router'
import Input from '@/components/Input.vue'
import Button from '@/components/Button.vue'
import RequestError from '@/components/RequestError.vue'
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { combineDateAndTime } from '@/utils/date'
import { createEvent as createEventApi } from '@/service/events'

const route = useRoute()
const categoryId = route.params.categoryId as string
const categoryName = route.query.categoryName as string

const eventToCreate = ref({
  title: categoryName || null,
  date: null,
  time: null,
  location: null,
  description: null
})

const createLoading = ref(false)
const validationErrors = ref<Record<string, string>>({})

const minDate = new Date().toISOString().split('T')[0]

const validate = () => {
  validationErrors.value = {}
  let isValid = true

  if (!eventToCreate.value.title) {
    validationErrors.value.title = i18n.global.t('events.new.required')
    isValid = false
  }
  if (!eventToCreate.value.date) {
    validationErrors.value.date = i18n.global.t('events.new.required')
    isValid = false
  }
  if (!eventToCreate.value.time) {
    validationErrors.value.time = i18n.global.t('events.new.required')
    isValid = false
  }
  if (!eventToCreate.value.location) {
    validationErrors.value.location = i18n.global.t('events.new.required')
    isValid = false
  }

  return isValid
}

watch(eventToCreate, () => {
  if (Object.keys(validationErrors.value).length > 0) {
    // Only clear errors if they exist to avoid unnecessary re-renders or logic
    if (eventToCreate.value.title) delete validationErrors.value.title
    if (eventToCreate.value.date) delete validationErrors.value.date
    if (eventToCreate.value.time) delete validationErrors.value.time
    if (eventToCreate.value.location) delete validationErrors.value.location
  }
}, { deep: true })

const createEvent = () => {
  if (!validate()) {
    return
  }

  createLoading.value = true;

  createEventApi({
    title: eventToCreate.value.title!,
    description: eventToCreate.value.description,
    startTime: combineDateAndTime(eventToCreate.value.date!, eventToCreate.value.time!),
    location: eventToCreate.value.location!,
    categoryId: categoryId
  })
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
  <div>
    <Header
      :title="i18n.global.t('events.new.header')"
      show-divider
      left-action="back"
    />

    <Input
      id="title"
      type="text"
      :placeholder="i18n.global.t('events.new.title')"
      v-model="eventToCreate.title"
      :error="validationErrors.title"
      required
    />
    <Input
      id="date"
      type="date"
      :placeholder="i18n.global.t('events.new.date')"
      v-model="eventToCreate.date"
      :error="validationErrors.date"
      required
      :min="minDate"
    />
    <Input
      id="time"
      type="time"
      :placeholder="i18n.global.t('events.new.time')"
      v-model="eventToCreate.time"
      :error="validationErrors.time"
      required
    />
    <Input
      id="location"
      type="text"
      :placeholder="i18n.global.t('events.new.location')"
      v-model="eventToCreate.location"
      :error="validationErrors.location"
      required
    />
    <Input
      id="description"
      type="text"
      :placeholder="i18n.global.t('events.new.description')"
      multi-line
      v-model="eventToCreate.description"
    />

    <RequestError class="newEventRequestError" />

    <Button
      :title="i18n.global.t('events.new.create')"
      :loading="createLoading"
      @click="createEvent"
      class="newEventButton"
    />
  </div>
</template>

<style scoped>
.newEventRequestError {
  margin-top: 24px;
}

.newEventButton {
  margin-top: 24px;
}
</style>

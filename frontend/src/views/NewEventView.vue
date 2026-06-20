<script setup lang="ts">
import Header from '@/components/Header.vue'
import i18n from '@/plugins/i18n'
import router from '@/router'
import Input from '@/components/Input.vue'
import Button from '@/components/Button.vue'
import RequestError from '@/components/RequestError.vue'
import { ref, watch } from 'vue'
import { useMediaQuery } from '@vueuse/core'

const isDesktop = useMediaQuery('(min-width: 768px)')
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
  } else if (isDesktop.value && !isValidTime(eventToCreate.value.time)) {
    validationErrors.value.time = i18n.global.t('events.new.timeInvalid')
    isValid = false
  }
  if (!eventToCreate.value.location) {
    validationErrors.value.location = i18n.global.t('events.new.required')
    isValid = false
  }

  return isValid
}

const isValidTime = (value: string) => /^([01]\d|2[0-3]):([0-5]\d)$/.test(value)

watch(eventToCreate, () => {
  if (Object.keys(validationErrors.value).length > 0) {
    if (eventToCreate.value.title) delete validationErrors.value.title
    if (eventToCreate.value.date) delete validationErrors.value.date
    if (eventToCreate.value.time && isValidTime(eventToCreate.value.time)) delete validationErrors.value.time
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
  <div class="newEventView">
    <div class="newEventCard">
      <Header
        :title="i18n.global.t('events.new.header')"
        show-divider
        left-action="back"
      />

      <div class="newEventForm">
        <Input
          id="title"
          type="text"
          :placeholder="i18n.global.t('events.new.title')"
          v-model="eventToCreate.title"
          :error="validationErrors.title"
          required
        />
        <div class="newEventDateTimeRow">
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
            :type="isDesktop ? 'text' : 'time'"
            :placeholder="isDesktop ? `${i18n.global.t('events.new.time')} (HH:MM)` : i18n.global.t('events.new.time')"
            v-model="eventToCreate.time"
            :error="validationErrors.time"
            required
          />
        </div>
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
    </div>
  </div>
</template>

<style scoped>
.newEventRequestError {
  margin-top: 24px;
}

.newEventButton {
  margin-top: 24px;
}

.newEventDateTimeRow {
  display: flex;
  flex-direction: column;
}

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

  .newEventDateTimeRow {
    flex-direction: row;
    gap: var(--gap);
  }

  .newEventForm {
    display: flex;
    flex-direction: column;
  }

  .newEventButton {
    width: auto;
    align-self: flex-end;
    padding-inline: calc(var(--padding) * 3);
  }
}
</style>

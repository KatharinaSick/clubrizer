<script setup lang="ts">
import { ref, watch } from 'vue'
import { useMediaQuery } from '@vueuse/core'
import i18n from '@/plugins/i18n'
import Input from '@/components/Input.vue'
import Button from '@/components/Button.vue'
import RequestError from '@/components/RequestError.vue'
import { combineDateAndTime, toDateInput, toTimeInput } from '@/utils/date'

export interface EventFormValues {
  title: string
  description: string | null
  startTime: Date
  location: string
}

// EventForm is the form for creating and editing an event. It validates the fields and
// emits the values on submit; saving them is up to the screen using it.
const props = defineProps<{
  submitTitle: string
  loading?: boolean
  // initialTitle prefills only the title, e.g. with the category name of a new event.
  initialTitle?: string
  initialValues?: EventFormValues
}>()

const emit = defineEmits<{
  submit: [values: EventFormValues]
}>()

const isDesktop = useMediaQuery('(min-width: 768px)')

const initial = props.initialValues
const form = ref({
  title: initial?.title ?? props.initialTitle ?? null,
  date: initial ? toDateInput(new Date(initial.startTime)) : null,
  time: initial ? toTimeInput(new Date(initial.startTime)) : null,
  location: initial?.location ?? null,
  description: initial?.description ?? null,
})

const validationErrors = ref<Record<string, string>>({})

const minDate = toDateInput(new Date())

const isValidTime = (value: string) => /^([01]\d|2[0-3]):([0-5]\d)$/.test(value)

const validate = () => {
  validationErrors.value = {}
  let isValid = true

  if (!form.value.title) {
    validationErrors.value.title = i18n.global.t('eventForm.required')
    isValid = false
  }
  if (!form.value.date) {
    validationErrors.value.date = i18n.global.t('eventForm.required')
    isValid = false
  }
  if (!form.value.time) {
    validationErrors.value.time = i18n.global.t('eventForm.required')
    isValid = false
  } else if (isDesktop.value && !isValidTime(form.value.time)) {
    validationErrors.value.time = i18n.global.t('eventForm.timeInvalid')
    isValid = false
  }
  if (!form.value.location) {
    validationErrors.value.location = i18n.global.t('eventForm.required')
    isValid = false
  }

  return isValid
}

watch(form, () => {
  if (Object.keys(validationErrors.value).length > 0) {
    if (form.value.title) delete validationErrors.value.title
    if (form.value.date) delete validationErrors.value.date
    if (form.value.time && isValidTime(form.value.time)) delete validationErrors.value.time
    if (form.value.location) delete validationErrors.value.location
  }
}, { deep: true })

const submit = () => {
  if (props.loading || !validate()) {
    return
  }
  emit('submit', {
    title: form.value.title!,
    description: form.value.description,
    startTime: combineDateAndTime(form.value.date!, form.value.time!),
    location: form.value.location!,
  })
}
</script>

<template>
  <div class="eventForm">
    <Input
      id="title"
      type="text"
      :placeholder="i18n.global.t('eventForm.title')"
      v-model="form.title"
      :error="validationErrors.title"
      required
    />
    <div class="eventFormDateTimeRow">
      <Input
        id="date"
        type="date"
        :placeholder="i18n.global.t('eventForm.date')"
        v-model="form.date"
        :error="validationErrors.date"
        required
        :min="minDate"
      />
      <Input
        id="time"
        :type="isDesktop ? 'text' : 'time'"
        :placeholder="isDesktop ? `${i18n.global.t('eventForm.time')} (HH:MM)` : i18n.global.t('eventForm.time')"
        v-model="form.time"
        :error="validationErrors.time"
        required
      />
    </div>
    <Input
      id="location"
      type="text"
      :placeholder="i18n.global.t('eventForm.location')"
      v-model="form.location"
      :error="validationErrors.location"
      required
    />
    <Input
      id="description"
      type="text"
      :placeholder="i18n.global.t('eventForm.description')"
      multi-line
      v-model="form.description"
    />

    <RequestError class="eventFormRequestError" />

    <Button
      :title="submitTitle"
      :loading="loading"
      @click="submit"
      class="eventFormButton"
    />
  </div>
</template>

<style scoped>
.eventFormRequestError {
  margin-top: calc(var(--padding) * 2);
}

.eventFormButton {
  margin-top: calc(var(--padding) * 2);
}

.eventFormDateTimeRow {
  display: flex;
  flex-direction: column;
}

@media (min-width: 768px) {
  .eventForm {
    display: flex;
    flex-direction: column;
  }

  .eventFormDateTimeRow {
    flex-direction: row;
    gap: var(--gap);
  }

  .eventFormButton {
    width: auto;
    align-self: flex-end;
    padding-inline: calc(var(--padding) * 3);
  }
}
</style>

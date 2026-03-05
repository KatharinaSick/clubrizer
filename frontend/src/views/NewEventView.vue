<script setup lang="ts">
import Header from '@/components/Header.vue'
import i18n from '@/plugins/i18n'
import router from '@/router'
import Input from '@/components/Input.vue'
import Button from '@/components/Button.vue'
import Alert from '@/components/Alert.vue'
import { ref } from 'vue'
import axios from '@/plugins/axios'
import { useRoute } from 'vue-router'

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
const createError = ref<string | null>(null)

const createEvent = () => {
  createLoading.value = true;
  createError.value = null;

  axios
    .post('/events', {
      title: eventToCreate.value.title,
      description: eventToCreate.value.description,
      startTime: new Date(eventToCreate.value.date + 'T' + eventToCreate.value.time),
      location: eventToCreate.value.location,
      categoryId: categoryId
    })
    .then(response => {
      createLoading.value = false;
      // Redirect to the new event
      // router.push(`/events/${response.data.id}`)
      // For now, redirect to list until detail view is ready
      router.push('/events')
    })
    .catch(error => {
      createLoading.value = false;
      createError.value = error.response?.data || i18n.global.t('events.new.failedToCreate');
    })
}
</script>

<template>
  <div>
    <Header
      :title="i18n.global.t('events.new.header')"
      show-divider
      left-action="back"
      @back="router.back()"
    />

    <Input
      id="title"
      type="text"
      :placeholder="i18n.global.t('events.new.title')"
      v-model="eventToCreate.title"
    />
    <Input
      id="date"
      type="date"
      :placeholder="i18n.global.t('events.new.date')"
      v-model="eventToCreate.date"
    />
    <Input
      id="time"
      type="time"
      :placeholder="i18n.global.t('events.new.time')"
      v-model="eventToCreate.time"
    />
    <Input
      id="location"
      type="text"
      :placeholder="i18n.global.t('events.new.location')"
      v-model="eventToCreate.location"
    />
    <Input
      id="description"
      type="text"
      :placeholder="i18n.global.t('events.new.description')"
      multi-line
      v-model="eventToCreate.description"
    />

    <Alert
      v-if="createError"
      :title="i18n.global.t('events.new.failedToCreate')"
      :message="createError"
      class="newEventAlert"
    />

    <Button
      :title="i18n.global.t('events.new.create')"
      :loading="createLoading"
      @click="createEvent"
    />
  </div>
</template>

<style scoped>
.newEventAlert {
  margin-top: 24px;
}
</style>

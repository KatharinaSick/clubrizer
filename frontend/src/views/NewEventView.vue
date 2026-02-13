<script setup lang="ts">
import Header from '@/components/Header.vue'
import i18n from '@/plugins/i18n'
import router from '@/router'
import Input from '@/components/Input.vue'
import Button from '@/components/Button.vue'
import Select, { type Option } from '@/components/Select.vue'
import { onMounted, ref } from 'vue'
import axios from '@/plugins/axios'

const eventToCreate = ref({
  title: null,
  date: null,
  time: null,
  location: null,
  category: null as Option | null,
  description: null
})

const createLoading = ref(false)
const createError = ref(null)

const categoriesLoading = ref(true)
const categoriesError = ref(null)
const categories = ref<Option[]>([])

onMounted(() => {
  categoriesLoading.value = true
  categoriesError.value = null

  axios
    .get('/events/categories')
    .then(response => {
      categoriesLoading.value = false
      categories.value = response.data.map((c: any) => ({ id: c.id, name: c.name }))
    })
    .catch(error => {
      // TODO handle
      console.log(error)
    })
})

const createEvent = () => {
  createLoading.value = true;
  createError.value = null;

  axios
    .post('/events', {
      title: eventToCreate.value.title,
      description: eventToCreate.value.description,
      startTime: new Date(eventToCreate.value.date + 'T' + eventToCreate.value.time),
      location: eventToCreate.value.location,
      category: eventToCreate.value.category?.id
    })
    .then(response => {
      createLoading.value = false;
      // TODO show success or go back
    })
    .catch(error => {
      createLoading.value = false;
      // TODO handle properly
      createError.value = error.response.data;
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
    <Select
      id="category"
      :placeholder="i18n.global.t('events.new.category')"
      :options="categories"
      :loading="categoriesLoading"
      v-model="eventToCreate.category"
    />
    <Input
      id="description"
      type="text"
      :placeholder="i18n.global.t('events.new.description')"
      multi-line
      v-model="eventToCreate.description"
    />

    <Button :title="i18n.global.t('events.new.create')" @click="createEvent" />
  </div>
</template>

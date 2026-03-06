<script setup lang="ts">
import Header from '@/components/Header.vue'
import axios from '@/plugins/axios'
import i18n from '@/plugins/i18n'
import FloatingActionButton, { type Action } from '@/components/FloatingActionButton.vue'
import { onMounted, ref } from 'vue'
import type { EventProps } from '@/components/Event.vue'
import Event from '@/components/Event.vue'
import { useRouter } from 'vue-router'
import Alert from '@/components/Alert.vue'

const router = useRouter()

const eventsLoading = ref(true)
const eventsError = ref(null)
const events = ref<EventProps[]>([])

const categoriesLoading = ref(false)
const categories = ref<Action[]>([])

onMounted(() => {
  loadEvents()
  loadCategories()
})

const loadEvents = () => {
  // TODO actually show a loading indicator
  eventsLoading.value = true
  eventsError.value = null

  axios
    .get('/events')
    .then(response => {
      eventsLoading.value = false
      events.value = response.data
    })
    .catch(error => {
      // TODO handle
      console.log(error)
    })
}

const loadCategories = () => {
  categoriesLoading.value = true
  axios
    .get('/events/categories')
    .then(response => {
      categoriesLoading.value = false
      categories.value = response.data.map((c: any) => ({
        id: c.id,
        label: c.name,
        color: c.color,
        onClick: () => router.push({
          name: 'new-event',
          params: { categoryId: c.id },
          query: { categoryName: c.name }
        })
      }))
    })
    .catch(error => {
      categoriesLoading.value = false
      // TODO handle
      console.log('Failed to load categories', error)
    })
}
</script>

<template>
  <div>
    <Header :title="i18n.global.t('events.header')" />
    <Alert
      :title="i18n.global.t('events.development.title')"
      :message="i18n.global.t('events.development.message')"
      variant="warning"
      style="margin-bottom: 12px;"
    />

    <Event v-for="event in events" :event="event" />

    <FloatingActionButton :actions="categories" :loading="categoriesLoading" />
  </div>
</template>

<script setup lang="ts">
import Header from '@/components/Header.vue'
import axios from '@/plugins/axios'
import i18n from '@/plugins/i18n'
import FloatingActionButton, { type Action } from '@/components/FloatingActionButton.vue'
import { computed, onMounted, ref } from 'vue'
import type { EventProps } from '@/components/Event.vue'
import Event from '@/components/Event.vue'
import { useRouter } from 'vue-router'
import Alert from '@/components/Alert.vue'
import Loader from '@/components/Loader.vue'

const router = useRouter()

const eventsLoading = ref(true)
const eventsError = ref<string | null>(null)
const events = ref<EventProps[]>([])

const categoriesLoading = ref(false)
const categoriesError = ref<string | null>(null)
const categories = ref<Action[]>([])

onMounted(() => {
  loadEvents()
  loadCategories()
})

const loadEvents = () => {
  eventsLoading.value = true
  eventsError.value = null

  axios
    .get('/events')
    .then(response => {
      eventsLoading.value = false
      events.value = response.data
    })
    .catch(error => {
      eventsLoading.value = false
      eventsError.value = error.response?.data || error.message || 'Unknown error'
    })
}

const loadCategories = () => {
  categoriesLoading.value = true
  categoriesError.value = null
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
      categoriesError.value = error.response?.data || error.message || 'Unknown error'
    })
}

const errorState = computed(() => {
  if (eventsError.value && categoriesError.value) {
    return {
      title: i18n.global.t('events.failedToLoad'),
      message: `${eventsError.value} | ${categoriesError.value}`
    }
  }
  if (eventsError.value) {
    return {
      title: i18n.global.t('events.failedToLoad'),
      message: eventsError.value
    }
  }
  if (categoriesError.value) {
    return {
      title: i18n.global.t('events.failedToLoadCategories'),
      message: categoriesError.value
    }
  }
  return null
})

const openEvent = (id: string) => {
  router.push({ name: 'event-detail', params: { id } })
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

    <Loader v-if="eventsLoading" />

    <Alert
      v-if="errorState"
      :title="errorState.title"
      :message="errorState.message"
      variant="error"
      style="margin-bottom: 12px;"
    />

    <div v-if="!eventsLoading && !eventsError">
      <Event
        v-for="event in events"
        :key="event.id"
        :event="event"
        @click="openEvent(event.id)"
        style="cursor: pointer;"
      />
    </div>

    <FloatingActionButton :actions="categories" :loading="categoriesLoading" />
  </div>
</template>

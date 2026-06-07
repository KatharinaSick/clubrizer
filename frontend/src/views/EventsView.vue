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
import RequestError from '@/components/RequestError.vue'

const router = useRouter()

const events = ref<EventProps[]>([])
const eventsLoaded = ref(false)
const categoriesLoading = ref(false)
const categories = ref<Action[]>([])

onMounted(() => {
  loadEvents()
  loadCategories()
})

const loadEvents = () => {
  axios
    .get('/events')
    .then(response => {
      events.value = response.data ?? []
    })
    .catch(() => {})
    .finally(() => {
      eventsLoaded.value = true
    })
}

const loadCategories = () => {
  categoriesLoading.value = true
  axios
    .get('/events/categories')
    .then(response => {
      categoriesLoading.value = false
      categories.value = response.data
        .filter((c: any) => c.canCreate)
        .map((c: any) => ({
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
    .catch(() => {
      categoriesLoading.value = false
    })
}

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

    <RequestError style="margin-bottom: 12px;" />

    <Event
      v-for="event in events"
      :key="event.id"
      :event="event"
      @click="openEvent(event.id)"
      style="cursor: pointer;"
    />

    <p v-if="eventsLoaded && events.length === 0" class="eventsEmpty">
      {{ $t('events.noEvents') }}
    </p>

    <FloatingActionButton
      v-if="categoriesLoading || categories.length > 0"
      :actions="categories"
      :loading="categoriesLoading"
    />
  </div>
</template>

<style scoped>
.eventsEmpty {
  text-align: center;
  color: var(--text-gray);
  font-size: var(--font-size-small);
  margin-top: 48px;
}
</style>
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

interface EventCategory {
  id: string
  name: string
  color: string
  canCreate: boolean
}

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
      categories.value = (response.data as EventCategory[])
        .filter(c => c.canCreate)
        .map(c => ({
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

const eventYear = (event: EventProps) => new Date(event.startTime).getFullYear()

// Show a year separator only where the year changes from the previous event, so a list
// that stays within one year shows no divider at all.
const showYearDivider = (index: number) => {
  const current = events.value[index]
  const previous = events.value[index - 1]
  if (!current || !previous) return false
  return eventYear(current) !== eventYear(previous)
}
</script>

<template>
  <div class="eventsView">
    <Header :title="i18n.global.t('events.header')" />
    <Alert
      :title="i18n.global.t('events.development.title')"
      :message="i18n.global.t('events.development.message')"
      variant="warning"
      style="margin-bottom: 12px;"
    />

    <RequestError style="margin-bottom: 12px;" />

    <div class="eventsGrid">
      <template v-for="(event, index) in events" :key="event.id">
        <div v-if="showYearDivider(index)" class="eventsYearDivider">
          <span class="eventsYearDividerLabel">{{ eventYear(event) }}</span>
        </div>
        <Event
          :event="event"
          @click="openEvent(event.id)"
          style="cursor: pointer;"
        />
      </template>
    </div>

    <p v-if="eventsLoaded && events.length === 0" class="eventsEmpty">
      {{ $t('events.noEvents') }}
    </p>

    <FloatingActionButton
      v-if="categoriesLoading || categories.length > 0"
      :actions="categories"
      :loading="categoriesLoading"
      :label="$t('events.fab')"
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

.eventsYearDivider {
  display: flex;
  align-items: center;
  gap: var(--gap);
  margin: var(--gap) 0;
}

.eventsYearDivider::before,
.eventsYearDivider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--gray);
}

.eventsYearDividerLabel {
  color: var(--text-gray);
  font-size: var(--font-size-small);
}

@media (min-width: 768px) {
  .eventsView {
    max-width: var(--content-max-width);
    margin: 0 auto;
  }

  .eventsGrid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: var(--gap);
  }

  .eventsYearDivider {
    grid-column: 1 / -1;
  }
}
</style>
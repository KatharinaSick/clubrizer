<script setup lang="ts">
import Header from '@/components/Header.vue'
import axios from '@/plugins/axios'
import i18n from '@/plugins/i18n'
import FloatingActionButton, { type Action } from '@/components/FloatingActionButton.vue'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import type { EventProps } from '@/components/Event.vue'
import Event from '@/components/Event.vue'
import { useRouter } from 'vue-router'
import RequestError from '@/components/RequestError.vue'
import IconArrowDown from '@/components/icons/IconArrowDown.vue'

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
const fabActions = ref<Action[]>([])
const selectedCategoryId = ref<string | null>(null)
const isFilterOpen = ref(false)

// Unique categories that actually appear in the loaded event list, preserving sort order.
const filterCategories = computed(() => {
  const seen = new Set<string>()
  const result: { id: string; name: string; color: string }[] = []
  for (const e of events.value) {
    if (!seen.has(e.category.id)) {
      seen.add(e.category.id)
      result.push({ id: e.category.id, name: e.category.name, color: e.category.color })
    }
  }
  return result
})

const showFilter = computed(() => filterCategories.value.length > 1)

const selectedLabel = computed(() =>
  selectedCategoryId.value === null
    ? i18n.global.t('events.filter.allEvents')
    : (filterCategories.value.find(c => c.id === selectedCategoryId.value)?.name ?? i18n.global.t('events.filter.allEvents'))
)

const filteredEvents = computed(() =>
  selectedCategoryId.value === null
    ? events.value
    : events.value.filter(e => e.category.id === selectedCategoryId.value),
)

const selectCategory = (id: string | null) => {
  selectedCategoryId.value = id
  isFilterOpen.value = false
}

const handleClickOutside = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.eventsFilter')) {
    isFilterOpen.value = false
  }
}

onMounted(() => document.addEventListener('click', handleClickOutside, true))
onBeforeUnmount(() => document.removeEventListener('click', handleClickOutside, true))

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
      fabActions.value = (response.data as EventCategory[])
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
  const current = filteredEvents.value[index]
  const previous = filteredEvents.value[index - 1]
  if (!current || !previous) return false
  return eventYear(current) !== eventYear(previous)
}
</script>

<template>
  <div class="eventsView">
    <Header :title="i18n.global.t('events.header')" />

    <div v-if="showFilter" class="eventsFilter">
      <button class="eventsFilterTrigger" @click.stop="isFilterOpen = !isFilterOpen">
        {{ selectedLabel }}
        <IconArrowDown
          class="eventsFilterArrow"
          :class="{ eventsFilterArrowOpen: isFilterOpen }"
        />
      </button>
      <div v-if="isFilterOpen" class="eventsFilterDropdown">
        <button
          class="eventsFilterOption"
          :class="{ eventsFilterOptionActive: selectedCategoryId === null }"
          @click="selectCategory(null)"
        >
          {{ $t('events.filter.allEvents') }}
        </button>
        <button
          v-for="cat in filterCategories"
          :key="cat.id"
          class="eventsFilterOption"
          :class="{ eventsFilterOptionActive: selectedCategoryId === cat.id }"
          @click="selectCategory(cat.id)"
        >
          <span class="eventsFilterOptionDot" :style="{ background: cat.color }" />
          {{ cat.name }}
        </button>
      </div>
    </div>
    <RequestError style="margin-bottom: 12px;" />

    <div class="eventsGrid">
      <template v-for="(event, index) in filteredEvents" :key="event.id">
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

    <p v-if="eventsLoaded && filteredEvents.length === 0" class="eventsEmpty">
      {{ $t('events.noEvents') }}
    </p>

    <FloatingActionButton
      v-if="categoriesLoading || fabActions.length > 0"
      :actions="fabActions"
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

.eventsFilter {
  position: relative;
  display: flex;
  justify-content: flex-end;
  margin-bottom: var(--gap);
}

.eventsFilterTrigger {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: none;
  border: none;
  padding: 0;
  cursor: pointer;
  color: var(--text-gray);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
  font-family: inherit;
  white-space: nowrap;
}

.eventsFilterArrow {
  width: 12px;
  height: 12px;
  flex-shrink: 0;
  transition: transform 0.15s ease;
}

.eventsFilterArrowOpen {
  transform: rotate(180deg);
}

.eventsFilterDropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 160px;
  background: var(--background-color);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  border: 1px solid var(--gray);
  overflow: hidden;
  z-index: 100;
}

.eventsFilterOption {
  display: flex;
  align-items: center;
  gap: var(--gap);
  width: 100%;
  padding: 10px var(--padding);
  background: none;
  border: none;
  text-align: left;
  font-size: var(--font-size-small);
  font-family: inherit;
  cursor: pointer;
  color: var(--text-color);
}

.eventsFilterOption:hover {
  background: var(--light-gray);
}

.eventsFilterOptionActive {
  font-weight: var(--font-weight-medium);
  color: var(--blue);
}

.eventsFilterOptionDot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
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
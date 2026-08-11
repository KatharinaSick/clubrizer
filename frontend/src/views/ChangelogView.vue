<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Header from '@/components/Header.vue'
import i18n from '@/plugins/i18n'

const { tm } = useI18n()

interface ChangelogEntry {
  date: string
  items: string[]
}

const entries = computed(() => tm('changelog.entries') as ChangelogEntry[])
</script>

<template>
  <div class="changelogView">
    <div class="changelogCard">
      <Header :title="i18n.global.t('changelog.header')" left-action="back" />
      <div class="changelogContent">
        <div v-for="(entry, index) in entries" :key="index" class="changelogEntry">
          <p class="changelogDate">{{ entry.date }}</p>
          <ul class="changelogList">
            <li v-for="(item, itemIndex) in entry.items" :key="itemIndex" class="changelogItem">
              {{ item }}
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.changelogContent {
  display: flex;
  flex-direction: column;
  gap: calc(var(--gap) * 3);
  padding: var(--padding) 0;
}

.changelogEntry {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.changelogDate {
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
  color: var(--text-gray);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.changelogList {
  margin: 0;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.changelogItem {
  line-height: 1.5;
}

@media (min-width: 768px) {
  .changelogView {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .changelogCard {
    width: 100%;
    max-width: var(--content-max-width);
    background: var(--background-color);
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
    padding: var(--padding);
    box-sizing: border-box;
  }
}
</style>

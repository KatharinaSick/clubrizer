<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router'
import Navigation from '@/components/Navigation.vue'
import TopProgressBar from '@/components/TopProgressBar.vue'
import { useRequestStore } from '@/stores/request'

const requestStore = useRequestStore()
useRouter().afterEach(() => requestStore.clearError())
</script>

<template>
  <div class="app">
    <TopProgressBar />
    <div class="content" :class="{ contentFullBleed: $route.meta.fullBleed }">
      <RouterView />
    </div>
    <Navigation v-if="$route.meta.showNavigation" class="navigation" />
  </div>
</template>

<style scoped>
.app {
  height: 100dvh;
  width: 100%;

  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content {
  flex-grow: 1;
  overflow-y: auto;
  padding: var(--padding);
  padding-bottom: var(--padding);
}

.contentFullBleed {
  padding: 0;
}

.navigation {
  flex-shrink: 0;
  margin-top: auto;

  margin-left: var(--padding);
  margin-right: var(--padding);
  margin-bottom: var(--padding);

  width: calc(100% - 2 * var(--padding));
}
</style>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import Navigation from '@/components/Navigation.vue'
</script>

<template>
  <div class="app">
    <div class="content">
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

.navigation {
  flex-shrink: 0;
  margin-top: auto;

  margin-left: var(--padding);
  margin-right: var(--padding);
  margin-bottom: var(--padding);

  width: calc(100% - 2 * var(--padding));
}
</style>

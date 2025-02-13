<script setup lang="ts">
import { RouterView } from 'vue-router'
import Navigation from '@/components/Navigation.vue'
</script>

<template>
  <div class="app">
    <RouterView class="router" />
    <Navigation v-if="$route.meta.showNavigation" />
  </div>
</template>

<style scoped>
.app {
  height: 100%;
  width: 100%;
  padding: var(--padding);

  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.router {
  flex-grow: 1;
  height: 100%;
}
</style>

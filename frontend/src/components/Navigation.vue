<script setup lang="ts">

import IconHome from '@/components/icons/IconHome.vue'
import { RouterLink } from 'vue-router'
import IconProfile from '@/components/icons/IconProfile.vue'
import IconEvents from '@/components/icons/IconEvents.vue'
</script>

<template>
  <nav class="navigation">
    <RouterLink to="/events" class="navigationLink" activeClass="navigationActiveLink">
      <IconEvents class="navigationIcon" />
    </RouterLink>
    <RouterLink to="/" class="navigationLink" activeClass="navigationActiveLink">
      <IconHome class="navigationIcon" />
    </RouterLink>
    <RouterLink to="/profile" class="navigationLink" activeClass="navigationActiveLink">
      <IconProfile class="navigationIcon" />
    </RouterLink>
  </nav>
</template>

<style scoped>
.navigation {
  width: 100%;

  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-evenly;

  box-shadow: var(--box-shadow);
  border-radius: var(--border-radius);
}

.navigationIcon {
  display: block;
  width: 48px;
  height: 28px;
  padding: 12px 12px;
}

.navigationLink {
  padding: 0;
  margin: 0;
  color: var(--gray);
}

.navigationActiveLink {
  color: var(--blue);
}
</style>

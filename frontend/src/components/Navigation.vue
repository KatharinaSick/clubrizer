<script setup lang="ts">

import { RouterLink, useRoute } from 'vue-router'

const route = useRoute()
import IconProfile from '@/components/icons/IconProfile.vue'
import IconEvents from '@/components/icons/IconEvents.vue'
</script>

<template>
  <nav class="navigation">
    <span class="navigationBrand">{{ $t('team') }}</span>
    <RouterLink to="/events" :class="['navigationLink', { navigationActiveLink: route.meta.activeNav === 'events' }]">
      <IconEvents class="navigationIcon" />
      <span class="navigationLabel">{{ $t('navigation.events') }}</span>
    </RouterLink>
    <RouterLink to="/profile" :class="['navigationLink', { navigationActiveLink: route.meta.activeNav === 'profile' }]">
      <IconProfile class="navigationIcon" />
      <span class="navigationLabel">{{ $t('navigation.profile') }}</span>
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

.navigationBrand {
  display: none;
}

.navigationLabel {
  display: none;
}

@media (min-width: 768px) {
  .navigation {
    justify-content: flex-end;
    border-radius: 0;
    box-sizing: border-box;
    padding: 0 max(var(--padding), calc((100% - var(--content-max-width)) / 2));
    z-index: 10;
  }

  .navigationBrand {
    display: block;
    margin-right: auto;
    font-weight: var(--font-weight-bold);
    background-image: var(--gradient);
    color: transparent;
    background-clip: text;
    -webkit-background-clip: text;
  }

  .navigationIcon {
    display: none;
  }

  .navigationLink {
    padding: 16px var(--padding);
    font-size: var(--font-size-medium);
    font-weight: var(--font-weight-medium);
    border-bottom: 2px solid transparent;
    text-decoration: none;
    color: var(--text-color);
  }

  .navigationActiveLink {
    color: var(--blue);
    border-bottom-color: var(--blue);
  }

  .navigationLabel {
    display: inline;
  }
}
</style>

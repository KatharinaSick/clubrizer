<script setup lang="ts">
import { useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import KidsManager from '@/components/KidsManager.vue'

const router = useRouter()
</script>

<template>
  <div class="manageKidsView">
    <div class="manageKidsCard">
      <Header
        :title="$t('manageKids.header')"
        left-action="back"
        :back-fn="() => router.push('/profile')"
        show-divider
      />

      <KidsManager />
    </div>
  </div>
</template>

<style scoped>
@media (min-width: 768px) {
  .manageKidsView {
    padding: var(--padding);
    box-sizing: border-box;
    display: flex;
    justify-content: center;
  }

  .manageKidsCard {
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

<script setup lang="ts">
import { useRequestStore } from '@/stores/request'

const requestStore = useRequestStore()
</script>

<template>
  <div class="topProgressBarTrack">
    <div v-if="requestStore.isLoading" class="topProgressBarFill" />
  </div>
</template>

<style scoped>
.topProgressBarTrack {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 3px;
  z-index: 9999;
}

.topProgressBarFill {
  height: 100%;
  background: var(--horizontal-gradient);
  animation: topProgressBarSlide 1.2s ease-in-out infinite;
}

@keyframes topProgressBarSlide {
  0% {
    width: 0%;
    margin-left: 0%;
  }
  50% {
    width: 60%;
    margin-left: 20%;
  }
  100% {
    width: 0%;
    margin-left: 100%;
  }
}
</style>
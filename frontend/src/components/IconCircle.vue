<script setup lang="ts">
// A white icon on a round gradient background — used to give a screen or card a friendly
// visual anchor (e.g. the bell on the notification screens).
</script>

<template>
  <span class="iconCircle" aria-hidden="true">
    <slot />
  </span>
</template>

<style scoped>
.iconCircle {
  flex-shrink: 0;
  width: 56px;
  height: 56px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--gradient);
  color: var(--white);
}

.iconCircle :deep(svg) {
  width: 28px;
  height: 28px;
}
</style>

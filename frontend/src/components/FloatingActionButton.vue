<script setup lang="ts">
import IconPlus from '@/components/icons/IconPlus.vue'
</script>

<template>
  <div class="fab">
    <IconPlus class="fabIcon"/>
  </div>
</template>

<style scoped>
.fab {
  width: 56px;
  height: 56px;

  position: absolute;
  right: var(--padding);
  bottom: calc(64px + var(--padding));

  display: flex;
  align-items: center;
  justify-content: center;

  background: var(--gradient);
  border-radius: 50%;
  box-shadow: var(--box-shadow);
}

.fabIcon {
  color: white;
  width: 36px;
  height: 36px;
}
</style>

<script setup lang="ts">
import IconBack from '@/components/icons/IconBack.vue'

defineProps<{
  title: string,
  showDivider?: boolean,
  leftAction?: 'back',
}>()
</script>

<template>
  <div class="header">
    <div class="headerTitleContainer">
      <div class="headerIcon">
        <IconBack class="headerIcon" v-if="leftAction === 'back'" @click="$emit('back')" />
      </div>
      <h1 class="headerTitle">{{ title }}</h1>
      <div class="headerIcon">
      </div>
    </div>
    <div v-if="showDivider" class="headerDivider" />
  </div>
</template>

<style scoped>
.header {
  width: 100%;
  text-align: center;
  padding: 4px 0;
}

.headerTitleContainer {
  width: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.headerTitle {
  margin: 0;
  padding: 0;
}

.headerIcon {
  width: 32px;
  height: 32px;
}

.headerDivider {
  height: 1px;
  width: 100%;
  background: var(--blue);
  opacity: 0.3;
}
</style>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import IconBack from '@/components/icons/IconBack.vue'
import Divider from '@/components/Divider.vue'

const props = defineProps<{
  title?: string,
  showDivider?: boolean,
  leftAction?: 'back',
  backFn?: () => void,
}>()

const router = useRouter()
</script>

<template>
  <div class="header">
    <div class="headerTitleContainer">
      <div class="headerIcon">
        <IconBack class="headerIcon" v-if="leftAction === 'back'" @click="props.backFn ? props.backFn() : router.back()" />
      </div>
      <h1 class="headerTitle">{{ title }}</h1>
      <div class="headerIcon">
      </div>
    </div>
    <Divider v-if="showDivider" />
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

</style>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  title?: string
  message: string
  variant?: 'error' | 'warning' | 'info'
  size?: 'small' | 'medium'
}>()

const variantClass = computed(() => {
  if (props.variant === 'warning') {
    return 'alertWarning'
  }
  if (props.variant === 'info') {
    return 'alertInfo'
  }
  return 'alertError'
})
</script>

<template>
  <div class="alert" :class="[variantClass, { alertSmall: size === 'small' }]">
    <p v-if="title" class="alertTitle">{{ title }}</p>
    <p class="alertMessage">{{ message }}</p>
    <div v-if="$slots.default" class="alertExtra">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.alert {
  width: 100%;
  padding: var(--padding);
  border-radius: var(--border-radius);
  box-sizing: border-box;
}

.alertSmall {
  font-size: var(--font-size-small);
}

.alertTitle {
  font-weight: var(--font-weight-bold);
}

.alertExtra {
  margin-top: var(--gap);
}

.alertError {
  background: var(--light-red);
  border: 1px solid var(--red);
}

.alertWarning {
  background: var(--light-orange);
  border: 1px solid var(--orange);
}

.alertInfo {
  background: var(--light-blue);
  border: 1px solid var(--blue);
}
</style>

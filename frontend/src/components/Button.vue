<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string
    loading?: boolean
    disabled?: boolean
    theme?: 'primary' | 'secondary'
  }>(),
  {
    theme: 'primary'
  }
)
</script>

<template>
  <button
    class="button"
    :class="theme"
    :disabled="disabled || loading"
  >
    <span v-if="loading">Loading...</span>
    <span v-else class="buttonContent">
      <span v-if="$slots.icon" class="buttonIconWrapper">
        <slot name="icon"></slot>
      </span>
      <span>{{ title }}</span>
    </span>
  </button>
</template>

<style scoped>
.button {
  width: 100%;
  padding: 12px;

  background: var(--horizotal-gradient);
  border-radius: var(--border-radius);
  border: 0;

  font-size: var(--font-size-medium);
  font-weight: var(--font-weight-medium);
  color: var(--white);
  cursor: pointer;

  display: flex;
  align-items: center;
  justify-content: center;
}

.buttonIconWrapper {
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.buttonContent {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--gap);
}

.button:active {
  background: var(--horizotal-gradient-active);
}

.button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background: var(--horizotal-gradient-active);
}

.button.secondary {
  background: var(--white);
  border: 1px solid var(--gray);
  color: var(--text-light);
  box-shadow: var(--box-shadow);
  border-radius: var(--gap);
  font-size: var(--font-size-small);
}

.button.secondary:active {
  background: var(--light-gray);
  border-color: var(--text-gray);
}

.button.secondary:disabled {
  background: var(--light-gray);
  border-color: var(--text-gray);
  opacity: 0.5;
}
</style>

<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string
    loading?: boolean
    disabled?: boolean
    theme?: 'primary' | 'secondary' | 'green' | 'red' | 'ghost' | 'tertiary'
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
    <span v-if="loading" class="buttonSpinner" />
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

  background: var(--horizontal-gradient);
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
  background: var(--horizontal-gradient-active);
}

.button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background: var(--horizontal-gradient-active);
}

.button.secondary {
  background: var(--white);
  border: 1px solid var(--gray);
  color: var(--text-light);
  box-shadow: var(--box-shadow);
  border-radius: var(--border-radius);
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

.button.green {
  background: var(--green);
  border: 1px solid var(--green);
  color: var(--white);
  box-shadow: var(--box-shadow);
}

.button.green:active {
  opacity: 0.9;
}

.button.green:disabled {
  opacity: 0.5;
  background: var(--green);
}

.button.red {
  background: var(--red);
  border: 1px solid var(--red);
  color: var(--white);
  box-shadow: var(--box-shadow);
}

.button.red:active {
  opacity: 0.9;
}

.button.red:disabled {
  opacity: 0.5;
  background: var(--red);
}

.button.ghost {
  background: var(--white);
  color: var(--blue);
}

.button.ghost:active {
  background: var(--light-gray);
}

.button.ghost:disabled {
  background: var(--light-gray);
  opacity: 0.7;
}

.button.tertiary {
  background: transparent;
  border: none;
  box-shadow: none;
  color: var(--text-gray);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-regular);
}

.button.tertiary:active {
  opacity: 0.7;
}

.button.tertiary:disabled {
  opacity: 0.4;
}

.buttonSpinner {
  display: inline-block;
  width: 18px;
  height: 18px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-right-color: currentColor;
  border-radius: 50%;
  animation: buttonSpin 0.7s linear infinite;
  flex-shrink: 0;
}

@keyframes buttonSpin {
  to { transform: rotate(360deg); }
}
</style>

<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string
    loading?: boolean
    disabled?: boolean
    theme?: 'primary' | 'secondary' | 'green' | 'red' | 'ghost' | 'tertiary'
    // iconOnly shows just the icon slot; the title is kept as the accessible label.
    iconOnly?: boolean
    // inline sizes the button to its content instead of full width.
    inline?: boolean
    // danger tints the text of the borderless themes (tertiary/ghost) red.
    danger?: boolean
    // accent tints the text of the borderless themes (tertiary/ghost) with the accent
    // colour, so a borderless button still reads as an actionable link rather than muted.
    accent?: boolean
  }>(),
  {
    theme: 'primary'
  }
)
</script>

<template>
  <button
    class="button"
    :class="[theme, { iconOnly, inline, danger, accent }]"
    :disabled="disabled || loading"
    :aria-label="iconOnly ? title : undefined"
  >
    <span v-if="loading" class="buttonSpinner" />
    <span v-else class="buttonContent">
      <span v-if="$slots.icon" class="buttonIconWrapper">
        <slot name="icon"></slot>
      </span>
      <span v-if="!iconOnly">{{ title }}</span>
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

.button.inline {
  width: auto;
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

.button.tertiary.danger,
.button.ghost.danger {
  color: var(--red);
}

.button.tertiary.accent,
.button.ghost.accent {
  color: var(--blue);
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

<script setup lang="ts">
defineProps<{
  id: string
  label: string
  description?: string
  disabled?: boolean
}>()

const value = defineModel<boolean>({ required: true })
</script>

<template>
  <label class="toggle" :class="{ toggleDisabled: disabled }" :for="id">
    <span class="toggleText">
      <span class="toggleLabel">{{ label }}</span>
      <span v-if="description" class="toggleDescription">{{ description }}</span>
    </span>
    <input
      :id="id"
      v-model="value"
      class="toggleInput"
      type="checkbox"
      role="switch"
      :disabled="disabled"
    />
    <span class="toggleTrack" aria-hidden="true">
      <span class="toggleThumb" />
    </span>
  </label>
</template>

<style scoped>
.toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--padding);
  cursor: pointer;
}

.toggleDisabled {
  cursor: default;
  opacity: 0.6;
}

.toggleText {
  display: flex;
  flex-direction: column;
}

.toggleDescription {
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

/* The real checkbox stays in the page for keyboards and screen readers, but is invisible. */
.toggleInput {
  position: absolute;
  opacity: 0;
  width: 1px;
  height: 1px;
}

.toggleTrack {
  flex-shrink: 0;
  position: relative;
  width: 44px;
  height: 24px;
  border-radius: 12px;
  background: var(--gray);
  transition: background 0.2s;
}

.toggleThumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--white);
  box-shadow: var(--box-shadow);
  transition: transform 0.2s;
}

.toggleInput:checked + .toggleTrack {
  background: var(--horizontal-gradient);
}

.toggleInput:checked + .toggleTrack .toggleThumb {
  transform: translateX(20px);
}

.toggleInput:focus-visible + .toggleTrack {
  outline: 2px solid var(--blue);
  outline-offset: 2px;
}
</style>

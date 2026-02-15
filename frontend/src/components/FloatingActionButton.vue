<script setup lang="ts">
import IconPlus from '@/components/icons/IconPlus.vue'
import { ref } from 'vue'

export interface Action {
  id: string
  label: string
  color?: string
  onClick: () => void
}

const props = defineProps<{
  actions?: Action[]
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'click'): void
}>()

const isOpen = ref(false)

const toggle = () => {
  if (props.loading) return

  if (props.actions && props.actions.length > 0) {
    isOpen.value = !isOpen.value
  } else {
    emit('click')
  }
}

const handleAction = (action: Action) => {
  action.onClick()
  isOpen.value = false
}
</script>

<template>
  <div class="fabContainer">
    <!-- Backdrop -->
    <div v-if="isOpen" class="fabBackdrop" @click="isOpen = false" />

    <!-- Actions List -->
    <div class="fabActions" :class="{ fabActionsOpen: isOpen }">
      <span
        v-for="action in actions"
        :key="action.id"
        class="fabActionLabel"
        :style="{ borderColor: action.color || 'var(--blue)' }"
        @click="handleAction(action)"
      >
        {{ action.label }}
      </span>
    </div>

    <!-- Main FAB -->
    <div
      class="fabButton"
      @click="toggle"
      :class="{ fabButtonOpen: isOpen, fabButtonLoading: loading }"
    >
      <div v-if="loading" class="fabSpinner"></div>
      <IconPlus v-else class="fabIcon" />
    </div>
  </div>
</template>

<style scoped>
.fabContainer {
  position: absolute;
  right: var(--padding);
  bottom: calc(64px + var(--padding));
  display: flex;
  flex-direction: column;
  align-items: end;
  z-index: 100;
}

.fabBackdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: var(--modal-background-color);
  z-index: 90;
}

.fabButton {
  width: 56px;
  height: 56px;

  display: flex;
  align-items: center;
  justify-content: center;

  background: var(--gradient);
  border-radius: 50%;
  box-shadow: var(--box-shadow);
  transition: transform 0.3s ease;
  z-index: 100;
  cursor: pointer;
}

.fabButtonOpen {
  transform: rotate(45deg);
}

.fabButtonLoading {
  cursor: not-allowed;
  opacity: 0.8;
}

.fabIcon {
  color: white;
  width: 36px;
  height: 36px;
}

.fabSpinner {
  width: 24px;
  height: 24px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s ease-in-out infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.fabActions {
  display: flex;
  flex-direction: column-reverse;
  align-items: end;
  gap: var(--padding);
  margin-bottom: var(--padding);

  opacity: 0;
  pointer-events: none;
  transform: translateY(24px);
  transition: all 0.3s ease;

  position: absolute;
  bottom: 100%;
  right: 0;
  z-index: 100;
}

.fabActionsOpen {
  opacity: 1;
  pointer-events: auto;
  transform: translateY(0);
}

.fabActionLabel {
  background: var(--background-color);
  padding: 10px 16px;
  border-radius: 24px;
  border: 2px solid transparent;

  font-weight: var(--font-weight-bold);
  color: var(--text-color);

  box-shadow: var(--box-shadow);
  white-space: nowrap;
  cursor: pointer;

  transition: transform 0.1s ease;
}

.fabActionLabel:active {
  transform: scale(0.95);
}
</style>

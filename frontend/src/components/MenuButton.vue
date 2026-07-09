<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import IconMore from '@/components/icons/IconMore.vue'

export interface MenuItem {
  label: string
  danger?: boolean
  onClick: () => void
}

defineProps<{
  items: MenuItem[]
  ariaLabel?: string
}>()

const open = ref(false)
const trigger = ref<HTMLButtonElement | null>(null)
const menuPosition = ref({ top: 0, right: 0 })

// The menu is teleported to <body>, so we position it from the trigger's
// on-screen rectangle: dropping below it and aligning their right edges.
function updatePosition() {
  const el = trigger.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  menuPosition.value = {
    top: rect.bottom + 8,
    right: window.innerWidth - rect.right,
  }
}

function openMenu() {
  updatePosition()
  open.value = true
}

function close() {
  open.value = false
}

function toggle() {
  if (open.value) {
    close()
  } else {
    openMenu()
  }
}

function select(item: MenuItem) {
  close()
  item.onClick()
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    close()
  }
}

// Scrolling or resizing would leave the menu detached from the trigger, so close it.
function onReflow() {
  if (open.value) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
  window.addEventListener('scroll', onReflow, true)
  window.addEventListener('resize', onReflow)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
  window.removeEventListener('scroll', onReflow, true)
  window.removeEventListener('resize', onReflow)
})
</script>

<template>
  <div class="menuButton">
    <button
      ref="trigger"
      type="button"
      class="menuButtonTrigger"
      :aria-label="ariaLabel"
      :aria-expanded="open"
      aria-haspopup="menu"
      @click="toggle"
    >
      <IconMore class="menuButtonIcon" />
    </button>

    <Teleport to="body">
      <div v-if="open" class="menuButtonScrim" @click="close">
        <div
          class="menuButtonMenu"
          role="menu"
          :style="{ top: `${menuPosition.top}px`, right: `${menuPosition.right}px` }"
          @click.stop
        >
          <button
            v-for="(item, index) in items"
            :key="index"
            type="button"
            role="menuitem"
            class="menuButtonItem"
            :class="{ menuButtonItemDanger: item.danger }"
            @click="select(item)"
          >
            {{ item.label }}
          </button>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.menuButton {
  display: inline-block;
}

.menuButtonTrigger {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(4px);
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--box-shadow);
  padding: 0;
}

.menuButtonIcon {
  width: 22px;
  height: 22px;
  color: var(--text-color);
}

.menuButtonScrim {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: var(--modal-background-color);
}

.menuButtonMenu {
  position: fixed;
  z-index: 1001;
  min-width: 180px;
  background: var(--background-color);
  border: 1px solid var(--gray);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  padding: var(--gap);
  display: flex;
  flex-direction: column;
}

.menuButtonItem {
  width: 100%;
  text-align: left;
  padding: var(--padding-input) var(--padding);
  background: transparent;
  border: none;
  border-radius: var(--border-radius);
  font-family: inherit;
  font-size: var(--font-size-medium);
  color: var(--text-color);
  cursor: pointer;
}

.menuButtonItem:hover,
.menuButtonItem:active {
  background: var(--light-gray);
}

.menuButtonItemDanger {
  color: var(--red);
}
</style>

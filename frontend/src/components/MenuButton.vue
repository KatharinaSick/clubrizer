<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import IconMore from '@/components/icons/IconMore.vue'
import Divider from '@/components/Divider.vue'

export interface MenuItem {
  label: string
  danger?: boolean
  // Shows a small red dot next to the item, e.g. to flag a section that needs attention.
  badge?: boolean
  // Renders the item small, gray and centered — for a low-emphasis link-style entry.
  subtle?: boolean
  // Draws a separator line above the item, to group it apart from the ones before it.
  divider?: boolean
  onClick: () => void
}

const props = withDefaults(defineProps<{
  items: MenuItem[]
  ariaLabel?: string
  // 'floating' (default) is the circular button meant to sit over content; 'bare' is a plain
  // icon button for use inside a header or toolbar.
  variant?: 'floating' | 'bare'
  // Shows a red dot on the trigger, e.g. when the menu contains something needing attention.
  badge?: boolean
  // Where the menu sits: 'below' (default) drops it just under the trigger; 'top-aligned' lines
  // the menu's top edge up with the trigger's top edge, so it grows downward from there.
  placement?: 'below' | 'top-aligned'
}>(), {
  variant: 'floating',
  badge: false,
  placement: 'below',
})

const open = ref(false)
const trigger = ref<HTMLButtonElement | null>(null)
const menuStyle = ref<Record<string, string>>({})

// The menu is teleported to <body>, so we position it from the trigger's on-screen
// rectangle: aligning their right edges, and anchoring above or below per `placement`.
function updatePosition() {
  const el = trigger.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const right = `${window.innerWidth - rect.right}px`
  menuStyle.value = props.placement === 'top-aligned'
    ? { top: `${rect.top}px`, right }
    : { top: `${rect.bottom + 8}px`, right }
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
      :class="variant === 'bare' ? 'menuButtonTriggerBare' : 'menuButtonTriggerFloating'"
      :aria-label="ariaLabel"
      :aria-expanded="open"
      aria-haspopup="menu"
      @click="toggle"
    >
      <slot name="icon"><IconMore class="menuButtonIcon" /></slot>
      <span v-if="badge" class="menuButtonBadge" aria-hidden="true" />
    </button>

    <Teleport to="body">
      <div v-if="open" class="menuButtonScrim" @click="close">
        <div
          class="menuButtonMenu"
          role="menu"
          :style="menuStyle"
          @click.stop
        >
          <template v-for="(item, index) in items" :key="index">
            <Divider v-if="item.divider" class="menuButtonDivider" />
            <button
              type="button"
              role="menuitem"
              class="menuButtonItem"
              :class="{ menuButtonItemDanger: item.danger, menuButtonItemSubtle: item.subtle }"
              @click="select(item)"
            >
              <span class="menuButtonItemLabel">{{ item.label }}</span>
              <span v-if="item.badge" class="menuButtonItemBadge" aria-hidden="true" />
            </button>
          </template>
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
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  position: relative;
  touch-action: manipulation;
}

.menuButtonTriggerFloating {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(4px);
  box-shadow: var(--box-shadow);
}

.menuButtonTriggerBare {
  background: transparent;
  color: var(--text-color);
}

.menuButtonTriggerBare:hover,
.menuButtonTriggerBare:active {
  background: var(--light-gray);
}

.menuButtonIcon {
  width: 22px;
  height: 22px;
  color: var(--text-color);
}

.menuButtonBadge {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: var(--blue);
  border: 2px solid var(--background-color);
  box-sizing: content-box;
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--gap);
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

.menuButtonItemSubtle {
  justify-content: center;
  font-size: var(--font-size-small);
  color: var(--text-gray);
}

.menuButtonDivider {
  margin: 4px 0;
}

.menuButtonItemBadge {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--blue);
  flex-shrink: 0;
}
</style>

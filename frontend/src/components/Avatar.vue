<script setup lang="ts">
import { computed, ref } from 'vue'

const SIZES = { sm: 32, md: 48, lg: 64, xl: 80 }

const props = withDefaults(defineProps<{
  givenName: string | null
  familyName: string | null
  nickName?: string | null
  picture?: string | null
  size?: keyof typeof SIZES
  gradient?: boolean
  label?: string
}>(), {
  size: 'md',
  gradient: false
})

const initials = computed(() => {
  return `${(props.givenName ?? '').charAt(0)}${(props.familyName ?? '').charAt(0)}`.toUpperCase()
})

const sizeInPx = computed(() => `${SIZES[props.size]}px`)

const labelVisible = ref(false)

const toggleLabel = () => {
  labelVisible.value = !labelVisible.value
}
</script>

<template>
  <div
    class="avatar"
    :class="{ avatarGradient: gradient, avatarClickable: !!label }"
    :style="{ '--avatarSize': sizeInPx }"
    @click="label && toggleLabel()"
  >
    <img v-if="picture" :src="picture" :alt="givenName ?? ''" class="avatarImage" />
    <div v-else class="avatarFallback">{{ initials }}</div>
    <span v-if="label" class="avatarLabel" :class="{ avatarLabelVisible: labelVisible }">
      {{ label }}
    </span>
  </div>
</template>

<style scoped>
.avatar {
  width: var(--avatarSize);
  height: var(--avatarSize);
  border-radius: 50%;
  flex-shrink: 0;
  position: relative;
}

.avatarClickable {
  cursor: pointer;
}

.avatarGradient {
  background: var(--gradient);
  padding: 1px;
}

.avatarImage {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.avatarFallback {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: var(--horizotal-gradient);
  color: var(--white);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
}

.avatarLabel {
  display: none;
  position: absolute;
  bottom: calc(100% + var(--gap));
  left: 0;
  background: var(--text-light);
  color: var(--white);
  font-size: var(--font-size-small);
  padding: calc(var(--gap) / 2) var(--gap);
  border-radius: var(--gap);
  white-space: nowrap;
  z-index: 1;
}

.avatarLabelVisible {
  display: block;
}
</style>

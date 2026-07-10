<script setup lang="ts">
import { computed } from 'vue'

const SIZES = { sm: 32, md: 48, lg: 64, xl: 80 }

const props = withDefaults(defineProps<{
  givenName: string | null
  familyName: string | null
  nickName?: string | null
  picture?: string | null
  size?: keyof typeof SIZES
  gradient?: boolean
}>(), {
  size: 'md',
  gradient: false
})

const initials = computed(() => {
  return `${(props.givenName ?? '').charAt(0)}${(props.familyName ?? '').charAt(0)}`.toUpperCase()
})

const sizeInPx = computed(() => `${SIZES[props.size]}px`)
</script>

<template>
  <div
    class="avatar"
    :class="{ avatarGradient: gradient }"
    :style="{ '--avatarSize': sizeInPx }"
  >
    <img v-if="picture" :src="picture" :alt="givenName ?? ''" class="avatarImage" />
    <div v-else class="avatarFallback">{{ initials }}</div>
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
  background: var(--horizontal-gradient);
  color: var(--white);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
}

</style>

<script setup lang="ts">
import { computed, defineAsyncComponent, ref } from 'vue'

// Loaded lazily to break the circular import: UserProfileModal renders an Avatar itself.
const UserProfileModal = defineAsyncComponent(() => import('@/components/UserProfileModal.vue'))

const SIZES = { sm: 32, md: 48, lg: 64, xl: 80 }

const props = withDefaults(defineProps<{
  givenName: string | null
  familyName: string | null
  nickName?: string | null
  picture?: string | null
  size?: keyof typeof SIZES
  gradient?: boolean
  // When true, the avatar becomes a button that opens the person's profile card on click.
  interactive?: boolean
  // For a kid's profile card: the display name of the managing parent.
  parent?: string | null
}>(), {
  size: 'md',
  gradient: false,
  interactive: false,
})

const initials = computed(() => {
  return `${(props.givenName ?? '').charAt(0)}${(props.familyName ?? '').charAt(0)}`.toUpperCase()
})

const sizeInPx = computed(() => `${SIZES[props.size]}px`)

const showProfile = ref(false)

function onClick() {
  if (props.interactive) showProfile.value = true
}
</script>

<template>
  <component
    :is="interactive ? 'button' : 'div'"
    :type="interactive ? 'button' : undefined"
    class="avatar"
    :class="{ avatarGradient: gradient, avatarInteractive: interactive }"
    :style="{ '--avatarSize': sizeInPx }"
    @click="onClick"
  >
    <img v-if="picture" :src="picture" :alt="givenName ?? ''" class="avatarImage" />
    <div v-else class="avatarFallback">{{ initials }}</div>

    <Teleport v-if="interactive" to="body">
      <UserProfileModal
        v-if="showProfile"
        :given-name="givenName ?? ''"
        :family-name="familyName ?? ''"
        :nick-name="nickName"
        :picture="picture"
        :parent="parent"
        @close="showProfile = false"
      />
    </Teleport>
  </component>
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

.avatarInteractive {
  display: block;
  padding: 0;
  border: none;
  background: none;
  cursor: pointer;
  appearance: none;
  font: inherit;
  color: inherit;
}

/* Keep the gradient ring even when the avatar is a clickable button. Higher specificity so
   it wins over the button reset above. */
.avatarGradient.avatarInteractive {
  padding: 1px;
  background: var(--gradient);
}

.avatarInteractive:hover {
  opacity: 0.85;
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

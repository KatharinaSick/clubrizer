<script setup lang="ts">
import Avatar from '@/components/Avatar.vue'

defineProps<{
  givenName: string
  familyName: string
  nickName?: string | null
  picture?: string | null
}>()

const emit = defineEmits<{ close: [] }>()
</script>

<template>
  <div class="userProfileModalOverlay" role="dialog" aria-modal="true" @click="emit('close')">
    <div class="userProfileModalBorder" @click.stop>
      <div class="userProfileModalContent">
        <button class="userProfileModalClose" aria-label="Close" @click="emit('close')">✕</button>
        <Avatar
          :picture="picture"
          :given-name="givenName"
          :family-name="familyName"
          size="xl"
          :gradient="true"
        />
        <div class="userProfileModalText">
          <strong v-if="nickName" class="userProfileModalNickname">{{ nickName }}</strong>
          <span class="userProfileModalFullName">{{ givenName }} {{ familyName }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.userProfileModalOverlay {
  position: fixed;
  inset: 0;
  z-index: 100;

  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--padding);
  box-sizing: border-box;

  background: var(--modal-background-color);
  animation: userProfileModalOverlayIn 0.15s ease-out;
}

.userProfileModalBorder {
  padding: 1px;
  background: var(--gradient);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  animation: userProfileModalCardIn 0.2s ease-out;
}

.userProfileModalContent {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--padding);
  padding: calc(var(--padding) * 2);
  background: var(--background-color);
  border-radius: calc(var(--border-radius) - 1px);
  min-width: 180px;
}

.userProfileModalClose {
  position: absolute;
  top: var(--gap);
  right: var(--gap);
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  background: var(--light-gray);
  color: var(--text-color);
  font-size: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

.userProfileModalText {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  text-align: center;
}

.userProfileModalNickname {
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.userProfileModalFullName {
  color: var(--text-color);
  font-size: var(--font-size-medium);
}

@keyframes userProfileModalOverlayIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes userProfileModalCardIn {
  from { opacity: 0; transform: translateY(12px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

@media (prefers-reduced-motion: reduce) {
  .userProfileModalOverlay,
  .userProfileModalBorder {
    animation: none;
  }
}
</style>

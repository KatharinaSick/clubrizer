<script setup lang="ts">
import type { User } from '@/stores/auth'

defineProps<{
  user: User
}>()
</script>

<template>
  <div class="profileInfo">
    <div class="profileInfoImageBorder">
      <img :src="user.picture" :alt="user.nickName" class="profileInfoImage" />
    </div>
    <div class="profileInfoText">
      <strong>{{ user.nickName }}</strong>
      <p>{{ user.givenName }} {{ user.familyName }}</p>
      <p>{{ user.email }}</p>
    </div>
  </div>
</template>

<style scoped>
.profileInfo {
  width: 100%;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: var(--padding);
  background-color: var(--white);
  margin: 0;
}

.profileInfoImageBorder {
  height: 80px;
  flex-basis: 80px;
  flex-grow: 0;
  flex-shrink: 0;
  border-radius: 50%;
  background: var(--gradient);
  padding: 1px;
}

.profileInfoImage {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.profileInfoText {
  width: 100%;
  overflow: hidden;
}
</style>

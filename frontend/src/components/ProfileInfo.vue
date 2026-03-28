<script setup lang="ts">
import type { User } from '@/stores/auth'
import Avatar from '@/components/Avatar.vue'

defineProps<{
  user: User
}>()
</script>

<template>
  <div class="profileInfo">
    <Avatar
      :picture="user.picture"
      :given-name="user.givenName"
      :family-name="user.familyName"
      :nick-name="user.nickName"
      size="xl"
      :gradient="true"
    />
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

.profileInfoText {
  width: 100%;
  overflow: hidden;
}
</style>

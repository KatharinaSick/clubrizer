<script setup lang="ts">
import type { User } from '@/stores/auth'
import type { Role } from '@/service/users'
import Avatar from '@/components/Avatar.vue'

defineProps<{
  user: User
  roles?: Role[]
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
      <strong class="profileInfoNickName">{{ user.nickName }}</strong>
      <p>{{ user.givenName }} {{ user.familyName }}</p>
      <p>{{ user.email }}</p>
      <div v-if="roles && roles.length > 0" class="profileInfoRoles">
        <span v-for="role in roles" :key="role.id" class="profileInfoRoleChip">
          {{ role.name }}
        </span>
      </div>
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

.profileInfoNickName {
  font-size: var(--font-size-large);
  font-weight: var(--font-weight-bold);
}

.profileInfoRoles {
  display: flex;
  flex-wrap: wrap;
  gap: var(--gap);
  margin-top: var(--gap);
}

.profileInfoRoleChip {
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-medium);
  background-color: var(--light-gray);
  color: var(--text-color);
  padding: 2px 10px;
  border-radius: 999px;
}
</style>

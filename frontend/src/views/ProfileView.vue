<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import ProfileInfo from '@/components/ProfileInfo.vue'
import Header from '@/components/Header.vue'
import Button from '@/components/Button.vue'
import i18n from '@/plugins/i18n'

const auth = useAuthStore()
const router = useRouter()

const logout = async () => {
  await auth.logout()
  router.push('/signin')
}
</script>

<template>
  <div>
    <Header :title="i18n.global.t('profile.header')" />
    <ProfileInfo :user="auth.user" />
    <Button :title="i18n.global.t('profile.logout')" @click="logout" />
  </div>
</template>

<style scoped>
</style>

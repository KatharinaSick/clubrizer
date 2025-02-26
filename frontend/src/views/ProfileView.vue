<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import ProfileInfo from '@/components/ProfileInfo.vue'
import Header from '@/components/Header.vue'
import i18n from '@/plugins/i18n'

const auth = useAuthStore()
</script>

<template>
  <div>
    <Header :title="i18n.global.t('profile.header')" />
    <ProfileInfo :user="auth.user" />
  </div>
</template>

<style scoped>
</style>

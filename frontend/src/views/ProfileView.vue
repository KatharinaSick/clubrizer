<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import ProfileInfo from '@/components/ProfileInfo.vue'
import Header from '@/components/Header.vue'

const auth = useAuthStore()
</script>

<template>
  <div>
    <Header title="My Profile" />
    <ProfileInfo :user="auth.user" />
  </div>
</template>

<style scoped>
</style>

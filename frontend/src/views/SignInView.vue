<script setup lang="ts">
import { type CredentialResponse, GoogleSignInButton } from 'vue3-google-signin'
import { useAuthStore } from '@/stores/auth'
import Loader from '@/components/Loader.vue'
import Modal from '@/components/Modal.vue'
import Alert from '@/components/Alert.vue'
import { ref } from 'vue'
import router from '@/router'
import { useRoute } from 'vue-router'

const auth = useAuthStore()
const route = useRoute()
const googleSignInError = ref<string>('')

const handleGoogleLoginSuccess = (response: CredentialResponse) => {
  const { credential } = response
  if (!credential) {
    googleSignInError.value = 'Failed to parse Google login token'
    return
  }

  auth.login(credential)
    .then(() => {
      if (auth.isLoggedIn) {
        router.replace({path: route.query.redirect as string ?? '/'})
      }
    })
}
</script>

<template>
  <div class="container">
    <div class="center">
      <img
        alt="Clubrizer Logo"
        class="logo"
        src="@/assets/logo.svg"
      />
      <h1 class="title">Clubrizer</h1>
      <h3 class="slogan">Team Up!</h3>
    </div>

    <div class="center">
      <h1>Welcome to<br />LISC-2010!</h1>
      <p>Sign in to get started</p>
    </div>

    <GoogleSignInButton
      @success="handleGoogleLoginSuccess"
      @error="() => googleSignInError = 'Failed to sign in with Google'"
      size="large"
      shape="pill"
    ></GoogleSignInButton>
    <Modal v-if="auth.isLoading">
      <Loader />
    </Modal>
    <Alert v-if="googleSignInError" class="alert" title="Failed to sign in" :message="googleSignInError"/>
    <Alert v-if="auth.error" class="alert" title="Failed to sign in" :message="auth.error"/>

  </div>
</template>

<style scoped>
.container {
  width: 100%;
  margin-top: 64px;

  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 48px;
}

.center {
  text-align: center;
}

.logo {
  height: 100px;
}

.title {
  margin: 8px 0 0 0;
  text-transform: uppercase;
  letter-spacing: .1rem;
}

.slogan {
  margin: 0;
  text-transform: uppercase;
  letter-spacing: .1rem;
  background-image: var(--gradient);
  color: transparent;
  background-clip: text;
}
</style>

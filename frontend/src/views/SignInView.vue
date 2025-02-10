<script setup lang="ts">
import {
  GoogleSignInButton,
  type CredentialResponse,
  decodeCredential
} from 'vue3-google-signin'
import { ref } from 'vue'
import axios from '@/plugins/axios'

const profileInfo = ref()

const handleGoogleLoginSuccess = (response: CredentialResponse) => {
  const { credential } = response;
  if (!credential) {
    // TODO error
    return
  }

  axios
    .post("/users", {idToken: credential})
    .then((response) => console.log(response))
    .catch((error) => console.error(error));

  /*const googleProfileInfo = decodeCredential(credential);
  profileInfo.value = {
    email: googleProfileInfo.email,
    fullName: googleProfileInfo.name,
    nickName: googleProfileInfo.given_name,
    picture: googleProfileInfo.picture,
  }*/
};

const handleGoogleLoginError = () => {
  // TODO do something
  console.error("Login failed");
};
</script>

<template>
  <div>
    <h1>Sign In</h1>

    <GoogleSignInButton
      v-if="!profileInfo"
      @success="handleGoogleLoginSuccess"
      @error="handleGoogleLoginError"
      size="large"
      shape="pill"
    ></GoogleSignInButton>

    <div v-if="profileInfo">
      <h2>Tell us a bit about yourself</h2>
      TODO allow user to customize information here with pre-filled input fields

      <strong>TODO just skip this step. People can edit their profile later on anyways.
      Still, enrich token with roll & approval information in the backend and return access & refresh token
        Create axios interceptor to automatically refresh token
      </strong>
      <button></button>
    </div>
  </div>
</template>

<style>

</style>

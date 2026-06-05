<script setup lang="ts">
// Sign-in is being replaced with OTP — placeholder view for now
</script>

<template>
  <div class="container">
    <div class="center">
      <img alt="Clubrizer Logo" class="logo" src="@/assets/logo.svg" />
      <h1 class="title">Clubrizer</h1>
      <h3 class="slogan">Team Up!</h3>
    </div>
    <div class="center">
      <h1>{{ $t('signIn.welcomeTo') }}<br />{{ $t('team') }}!</h1>
      <p>{{ $t('signIn.getStarted') }}</p>
    </div>
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


<script setup lang="ts">
import Header from '@/components/Header.vue'
import axios from '@/plugins/axios'

const sendRequest = () => {
  axios.get('/events').then(response => console.log(response))
}
</script>

<template>
  <div>
  <Header title="Events" />

  <button @click="sendRequest">Send request</button>
  </div>
</template>

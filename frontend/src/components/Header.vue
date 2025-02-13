<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <div class="header">
    <h1>{{title}}</h1>
  </div>
</template>

<style scoped>
.header {
  width: 100%;
  text-align: center;
  padding: 4px 0;
}
</style>

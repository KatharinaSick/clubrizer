<script setup lang="ts">
// A numbered list of instructions, each number in a gradient circle — e.g. how to add the
// app to the iPhone home screen. Takes the already translated steps, so every screen keeps
// its own wording.
defineProps<{
  steps: string[]
}>()
</script>

<template>
  <ol class="numberedSteps">
    <li v-for="(step, index) in steps" :key="index" class="numberedStepsStep">
      <span class="numberedStepsNumber">{{ index + 1 }}</span>
      <span class="numberedStepsText">{{ step }}</span>
    </li>
  </ol>
</template>

<style scoped>
.numberedSteps {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--padding);
  text-align: left;
}

.numberedStepsStep {
  display: flex;
  align-items: flex-start;
  gap: var(--padding);
}

.numberedStepsNumber {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--gradient);
  color: var(--white);
  font-size: var(--font-size-small);
  font-weight: var(--font-weight-bold);
}

/* Lines the first line of text up with the middle of the number. */
.numberedStepsText {
  padding-top: calc(var(--gap) / 2);
  font-size: var(--font-size-small);
  line-height: 1.5;
}
</style>

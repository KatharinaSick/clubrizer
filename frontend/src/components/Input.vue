<script setup lang="ts">
defineProps<{
  id: string
  type: 'text' | 'date' | 'time'

  placeholder?: string
  multiLine?: boolean
}>()

const value = defineModel()

</script>

<template>
  <div class="inputWrapper">
    <textarea
      v-if="multiLine"
      class="input"
      :id="id"
      placeholder=""
      rows="5"
      v-model="value"
      required
    ></textarea>
    <input
      v-else-if="type === 'date'"
      class="input"
      :id="id"
      placeholder=""
      type="text"
      onfocus="(this.type = 'date')"
      onblur="if(!this.value) this.type = 'text'"
      v-model="value"
      required
    />
    <input
      v-else-if="type === 'time'"
      class="input"
      :id="id"
      placeholder=""
      type="text"
      onfocus="(this.type = 'time')"
      onblur="if(!this.value) this.type = 'text'"
      v-model="value"
      required
    />
    <input
      v-else
      class="input"
      :id="id"
      placeholder=""
      :type="type"
      v-model="value"
      required
    />
    <label
      :for="id"
      class='inputPlaceholder'>
      {{ placeholder }}
    </label>
  </div>
</template>

<style scoped>
.inputWrapper {
  width: 100%;
  margin: 0;
  box-sizing: border-box;

  padding-top: 24px;

  position: relative;
  display: flex;
  align-items: start;
}

.input {
  width: 100%;
  background: var(--light-gray);
  border-radius: var(--border-radius);
  border: 0;

  padding: var(--input-padding);
  outline: none;
  font-size: var(--font-size-medium);
  font-family: inherit;
}

.inputPlaceholder {
  position: absolute;
  left: var(--padding);
  top: 34px;
  transition: all 0.2s ease-in;
  color: var(--text-gray);
}

.input:is(:focus, :valid) ~ .inputPlaceholder {
  transform: translatey(calc(-28px));
  font-size: var(--font-size-small);
  color: var(--blue);
  left: 0;
}
</style>

<script setup lang="ts">
defineProps<{
  id: string
  type: 'text' | 'email' | 'date' | 'time'

  placeholder?: string
  multiLine?: boolean
  error?: string
  required?: boolean
  min?: string
  inputMode?: 'text' | 'search' | 'none' | 'tel' | 'url' | 'email' | 'numeric' | 'decimal'
  maxLength?: number
  theme?: 'default' | 'ghost'
}>()

const value = defineModel<string | number | null>()

</script>

<template>
  <div class="inputWrapper" :class="{ inputWrapperGhost: theme === 'ghost' }">
    <textarea
      v-if="multiLine"
      class="input"
      :class="{ 'inputError': error }"
      :id="id"
      placeholder=""
      rows="5"
      v-model="value"
      required
    ></textarea>
    <input
      v-else-if="type === 'date'"
      class="input"
      :class="{ 'inputError': error }"
      :id="id"
      type="date"
      v-model="value"
      :min="min"
      required
    />
    <input
      v-else-if="type === 'time'"
      class="input"
      :class="{ 'inputError': error }"
      :id="id"
      type="time"
      v-model="value"
      required
    />
    <input
      v-else
      class="input"
      :class="{ 'inputError': error }"
      :id="id"
      placeholder=""
      :type="type"
      v-model="value"
      required
      :inputmode="inputMode"
      :maxlength="maxLength"
    />
    <label
      :for="id"
      class='inputPlaceholder'>
      {{ placeholder }}<span v-if="required">*</span>
    </label>
    <span v-if="error" class="errorMessage">{{ error }}</span>
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
  flex-direction: column;
  align-items: start;
}

.input {
  width: 100%;
  background: var(--light-gray);
  border-radius: var(--border-radius);
  border: 1px solid transparent;

  padding: var(--padding-input);
  outline: none;
  font-size: var(--font-size-medium);
  font-family: inherit;
  box-sizing: border-box;
}

.inputError {
  border-color: var(--red);
}

.inputPlaceholder {
  position: absolute;
  left: var(--padding);
  top: 34px;
  transition: all 0.2s ease-in;
  color: var(--text-gray);
  pointer-events: none;
}

.input:is(:focus, :valid) ~ .inputPlaceholder {
  transform: translatey(calc(-28px));
  font-size: var(--font-size-small);
  color: var(--blue);
  left: 0;
}

.input.inputError:is(:focus, :valid) ~ .inputPlaceholder {
  color: var(--red);
}

.input[type="date"]:not(:focus):not(:valid),
.input[type="time"]:not(:focus):not(:valid) {
  color: transparent;
}

.errorMessage {
  color: var(--red);
  font-size: var(--font-size-small);
  margin-top: 4px;
  margin-left: var(--padding);
}

.inputWrapperGhost .input {
  background: var(--white);
}

.inputWrapperGhost .input:is(:focus, :valid) ~ .inputPlaceholder {
  color: var(--white);
}

.inputWrapperGhost .input.inputError:is(:focus, :valid) ~ .inputPlaceholder {
  color: var(--red);
}
</style>

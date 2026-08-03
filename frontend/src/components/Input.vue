<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue'

const props = defineProps<{
  id: string
  type: 'text' | 'email' | 'date' | 'time'

  placeholder?: string
  multiLine?: boolean
  // autoGrow (multiLine only) starts the field at a single row and grows it with the
  // content up to a max height, then scrolls — a compact chat-style composer instead of
  // the fixed five-row box.
  autoGrow?: boolean
  error?: string
  required?: boolean
  min?: string
  inputMode?: 'text' | 'search' | 'none' | 'tel' | 'url' | 'email' | 'numeric' | 'decimal'
  maxLength?: number
  theme?: 'default' | 'ghost'
}>()

const value = defineModel<string | number | null>()

// Largest height (px) an auto-growing field expands to before it starts scrolling.
const AUTO_GROW_MAX_HEIGHT = 120

const textarea = ref<HTMLTextAreaElement | null>(null)

// Recompute the auto-grow height: shrink to nothing first so it can grow or shrink to
// exactly fit the content, capped at the max. Runs after the DOM updates so the new
// content is measured.
function resize() {
  const el = textarea.value
  if (!el || !props.autoGrow) return
  el.style.height = 'auto'
  el.style.height = `${Math.min(el.scrollHeight, AUTO_GROW_MAX_HEIGHT)}px`
}

onMounted(() => {
  if (props.autoGrow) nextTick(resize)
})

// Grow/shrink whenever the value changes — including when the parent clears it after a
// submit, which snaps the field back to a single row.
watch(value, () => {
  if (props.autoGrow) nextTick(resize)
})

// Trim leading/trailing whitespace when the field loses focus. Done on blur (not
// while typing) so spaces inside the text are kept; a value of only blanks becomes
// empty, so `required` checks treat it as empty.
function trimOnBlur() {
  if (typeof value.value === 'string') {
    value.value = value.value.trim()
  }
}

</script>

<template>
  <div class="inputWrapper" :class="{ inputWrapperGhost: theme === 'ghost' }">
    <textarea
      v-if="multiLine"
      ref="textarea"
      class="input"
      :class="{ 'inputError': error, 'inputAutoGrow': autoGrow }"
      :id="id"
      placeholder=""
      :rows="autoGrow ? 1 : 5"
      :maxlength="maxLength"
      v-model="value"
      required
      @blur="trimOnBlur"
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
      @blur="trimOnBlur"
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

.inputAutoGrow {
  resize: none;
  overflow-y: auto;
  max-height: 120px;
  /* Hide the scrollbar (Firefox) but keep the field scrollable once it hits max height. */
  scrollbar-width: none;
}

/* Hide the scrollbar (WebKit/Chromium) too. */
.inputAutoGrow::-webkit-scrollbar {
  display: none;
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

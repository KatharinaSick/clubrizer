<template>
  <div class="modal">
    <div class="modalGradient">
      <div class="modalContent">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;

  width: 100%;
  height: 100%;

  display: flex;
  align-items: center;
  justify-content: center;

  background: var(--modal-background-color);
}

.modalGradient {
  padding: 1px;
  background: var(--gradient);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
}

.modalContent {
  padding: var(--padding);
  background: var(--background-color);
  border-radius: var(--border-radius);
}
</style>

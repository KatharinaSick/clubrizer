<template>
  <div class="modalOverlay">
    <div class="modalCard">
      <div class="modalContent">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped>
.modalOverlay {
  position: fixed;
  inset: 0;
  z-index: 100;

  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--padding);
  box-sizing: border-box;

  background: var(--modal-background-color);
  animation: modalOverlayIn 0.15s ease-out;
}

.modalCard {
  width: 100%;
  max-width: 400px;

  padding: 1px;
  background: var(--gradient);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);

  animation: modalCardIn 0.2s ease-out;
}

.modalContent {
  padding: calc(var(--padding) * 2);
  background: var(--background-color);
  border-radius: calc(var(--border-radius) - 1px);
}

@keyframes modalOverlayIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes modalCardIn {
  from { opacity: 0; transform: translateY(12px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

@media (prefers-reduced-motion: reduce) {
  .modalOverlay,
  .modalCard {
    animation: none;
  }
}
</style>

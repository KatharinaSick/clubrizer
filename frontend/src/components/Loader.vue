<template>
  <div class="loader"></div>
</template>

<style scoped>
.loader {
  width: 40px;
  aspect-ratio: 2;
  --_blue: no-repeat radial-gradient(circle closest-side,var(--blue) 90%,#0000);
  --_light-blue: no-repeat radial-gradient(circle closest-side,var(--light-blue) 90%,#0000);
  --_green: no-repeat radial-gradient(circle closest-side,var(--green) 90%,#0000);

  background:
    var(--_light-blue) 0   50%,
    var(--_blue) 50%  50%,
    var(--_green) 100% 50%;
  background-size: calc(100%/3) 50%;
  animation: l3 1s infinite linear;
}
@keyframes l3 {
  20%{background-position: 0   0, 50%  50%,100%  50%}
  40%{background-position: 0 100%, 50%   0,100%  50%}
  60%{background-position: 0  50%, 50% 100%,100%   0}
  80%{background-position: 0  50%, 50%  50%,100% 100%}
}
</style>

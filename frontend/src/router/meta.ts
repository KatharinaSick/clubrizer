import 'vue-router'

// To ensure it is treated as a module, add at least one `export` statement
export {}

declare module 'vue-router' {
  interface RouteMeta {
    // This (annoyingly) uses show instead of hide because of https://github.com/vuejs/vue-router/issues/3518
    showNavigation?: boolean
    requiresAuth?: boolean
  }
}

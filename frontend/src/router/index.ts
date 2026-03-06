import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/events'
    },
    /*
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { showNavigation: true, requiresAuth: true }
    },
    */
    {
      path: '/events',
      name: 'events',
      component: () => import('../views/EventsView.vue'),
      meta: { showNavigation: true, requiresAuth: true }
    },
    {
      path: '/events/new/:categoryId',
      name: 'new-event',
      component: () => import('../views/NewEventView.vue'),
      meta: { showNavigation: false, requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { showNavigation: true, requiresAuth: true }
    },
    {
      path: '/signin',
      name: 'singin',
      component: () => import('../views/SignInView.vue')
    }
  ]
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    return {
      path: '/signin',
      query: { redirect: to.fullPath }
    }
  } else if (!to.meta.requiresAuth && auth.isLoggedIn) {
    return { path: '/events' }
  }
})

export default router

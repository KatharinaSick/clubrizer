import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/events'
    },
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
      path: '/events/:id',
      name: 'event-detail',
      component: () => import('../views/EventDetailView.vue'),
      meta: { showNavigation: true, requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { showNavigation: true, requiresAuth: true }
    },
    {
      path: '/signin',
      name: 'signin',
      component: () => import('../views/SignInView.vue')
    },
    {
      path: '/pending-approval',
      name: 'pending-approval',
      component: () => import('../views/PendingApprovalView.vue')
    },
    {
      path: '/profile-setup',
      name: 'profile-setup',
      component: () => import('../views/ProfileSetupView.vue')
    },
    {
      path: '/design-system',
      name: 'design-system',
      component: () => import('../views/DesignSystemView.vue'),
      meta: { showNavigation: false }
    },
  ]
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  const isApproved = auth.user.status === 'approved'
  const needsProfileSetup = isApproved && auth.user.givenName === null

  if (!auth.isLoggedIn) {
    if (to.meta.requiresAuth || to.path === '/pending-approval' || to.path === '/profile-setup') {
      return { path: '/signin', query: { redirect: to.fullPath } }
    }
    return
  }

  if (!isApproved) {
    if (to.path !== '/pending-approval') {
      return { path: '/pending-approval' }
    }
    return
  }

  if (needsProfileSetup) {
    if (to.path !== '/profile-setup') {
      return { path: '/profile-setup' }
    }
    return
  }

  // User is logged in, approved, and has a complete profile — redirect away from auth-only pages
  if (['/signin', '/pending-approval', '/profile-setup'].includes(to.path)) {
    return { path: '/events' }
  }
})

export default router

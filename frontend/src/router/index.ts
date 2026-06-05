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
      name: 'singin',
      component: () => import('../views/SignInView.vue')
    },
    {
      path: '/pending-approval',
      name: 'pending-approval',
      component: () => import('../views/PendingApprovalView.vue')
    },
    {
      path: '/design-system',
      name: 'design-system',
      component: () => import('../views/DesignSystemView.vue'),
      meta: { showNavigation: false }
    }
  ]
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  const isApproved = auth.user.status === 'approved'

  if (!auth.isLoggedIn) {
    if (to.meta.requiresAuth || to.path === '/pending-approval') {
      return { path: '/signin', query: { redirect: to.fullPath } }
    }
    return
  }

  if (!isApproved) {
    // Guard against infinite redirect — let the user land on /pending-approval
    if (to.path !== '/pending-approval') {
      return { path: '/pending-approval' }
    }
    return
  }

  // User is logged in and approved — redirect away from auth-only pages
  if (to.path === '/signin' || to.path === '/pending-approval') {
    return { path: '/events' }
  }
})

export default router

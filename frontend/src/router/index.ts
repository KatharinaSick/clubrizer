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
      meta: { showNavigation: true, activeNav: 'events', requiresAuth: true }
    },
    {
      path: '/events/new/:categoryId',
      name: 'new-event',
      component: () => import('../views/NewEventView.vue'),
      meta: { showNavigation: false, activeNav: 'events', requiresAuth: true }
    },
    {
      path: '/events/:id',
      name: 'event-detail',
      component: () => import('../views/EventDetailView.vue'),
      meta: { showNavigation: false, fullBleed: true, activeNav: 'events', requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { showNavigation: true, activeNav: 'profile', requiresAuth: true }
    },
    {
      path: '/my-kids',
      name: 'my-kids',
      component: () => import('../views/ManageKidsView.vue'),
      meta: { showNavigation: false, activeNav: 'profile', requiresAuth: true }
    },
    {
      path: '/signin',
      name: 'signin',
      component: () => import('../views/SignInView.vue'),
      meta: { hideNavigation: true }
    },
    {
      path: '/onboarding',
      name: 'onboarding',
      component: () => import('../views/OnboardingView.vue'),
      meta: { hideNavigation: true }
    },
    {
      path: '/pending-approval',
      name: 'pending-approval',
      component: () => import('../views/PendingApprovalView.vue'),
      meta: { hideNavigation: true }
    },
    {
      path: '/account-setup',
      name: 'account-setup',
      component: () => import('../views/AccountSetupView.vue'),
      meta: { hideNavigation: true }
    },
    {
      path: '/getting-started',
      name: 'getting-started',
      component: () => import('../views/GettingStartedView.vue'),
      meta: { showNavigation: false }
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
  const status = auth.user.status
  const isApproved = status === 'approved'
  // Every approved account — guardians included — must fill in at least a first and last
  // name once, on the account-setup screen. There's no stored "done" flag: the presence of
  // the own name is the signal, and account-setup saves it last, so setup can't be
  // half-finished. (Kids' last names are enforced within that screen.)
  const needsAccountSetup =
    isApproved && (auth.user.givenName === null || auth.user.familyName === null)

  // Paths that need a session; a logged-out visitor is bounced to /signin. Note this
  // must NOT include /signin itself, or a logged-out visitor would redirect to
  // /signin forever.
  const sessionPaths = ['/onboarding', '/pending-approval', '/account-setup']

  if (!auth.isLoggedIn) {
    if (to.meta.requiresAuth || sessionPaths.includes(to.path)) {
      return { path: '/signin', query: { redirect: to.fullPath } }
    }
    return
  }

  // Phase 1 of signup: choose account type + add kids (active). Distinct from waiting.
  if (status === 'onboarding') {
    if (to.path !== '/onboarding') {
      return { path: '/onboarding' }
    }
    return
  }

  // Phase 2: submitted, waiting for an admin (or rejected) — passive.
  if (status === 'pending' || status === 'rejected') {
    if (to.path !== '/pending-approval') {
      return { path: '/pending-approval' }
    }
    return
  }

  if (needsAccountSetup) {
    if (to.path !== '/account-setup') {
      return { path: '/account-setup' }
    }
    return
  }

  // Logged in, approved, complete profile — redirect away from the auth-only pages.
  if (to.path === '/signin' || sessionPaths.includes(to.path)) {
    return { path: '/events' }
  }
})

export default router

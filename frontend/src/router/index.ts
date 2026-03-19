import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/pages/LoginPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/pages/RegisterPage.vue'),
      meta: { guest: true },
    },
    {
      path: '/',
      name: 'diary',
      component: () => import('@/pages/DiaryListPage.vue'),
      meta: { auth: true },
    },
    {
      path: '/entry/new',
      name: 'newEntry',
      component: () => import('@/pages/DiaryEntryFormPage.vue'),
      meta: { auth: true },
    },
    {
      path: '/entry/:id/edit',
      name: 'editEntry',
      component: () => import('@/pages/DiaryEntryFormPage.vue'),
      meta: { auth: true },
    },
    {
      path: '/privacy',
      name: 'privacy',
      component: () => import('@/pages/PrivacyPage.vue'),
    },
    {
      path: '/terms',
      name: 'terms',
      component: () => import('@/pages/TermsPage.vue'),
    },
    {
      path: '/imprint',
      name: 'imprint',
      component: () => import('@/pages/ImprintPage.vue'),
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.auth && !auth.isAuthenticated) {
    return '/login'
  }

  if (to.meta.guest && auth.isAuthenticated) {
    return '/'
  }
})

export default router

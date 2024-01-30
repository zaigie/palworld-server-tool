import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomePage.vue')
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminSystem/AdminSystem.vue')
      // beforeEnter: async (to, from, next) => {
      //   const token = localStorage.getItem('token')
      //   if (!token) {
      //     next({ name: 'login' })
      //   } else {
      //     next()
      //   }
      // }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginPage.vue')
      // beforeEnter: async (to, from, next) => {
      //   const token = localStorage.getItem('token')
      //   if (!token) {
      //     next()
      //   } else {
      //     next({ name: from.name })
      //   }
      // }
    }
  ]
})

export default router

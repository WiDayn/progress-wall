import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'


const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    redirect: (to) => {
      const userStore = useUserStore()
      if (userStore.isLoggedIn) {
        return '/teams'
      }
      return '/login'
    }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/user/LoginView.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/user/RegisterView.vue')
  },
  {
    path: '/teams',
    name: 'TeamList',
    component: () => import('@/views/team/TeamListView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/teams/:teamId(\\d+)/projects',
    name: 'ProjectList',
    component: () => import('@/views/project/ProjectListView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:projectId(\\d+)/boards',
    name: 'ProjectBoards',
    component: () => import('@/views/dashboard/BoardListView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/kanban/:boardId(\\d+)',
    name: 'Kanban',
    component: () => import('@/views/dashboard/KanbanView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/tasks/:taskId',
    name: 'TaskDetail',
    component: () => import('@/views/dashboard/TaskDetailView.vue'),
    props: true
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/user/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/user/SettingsView.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 初始化用户状态
  //userStore.initializeFromStorage()
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    // 需要登录但未登录，重定向到登录页
    next({
      path: '/login',
      query: { redirect: to.fullPath }
    })
  } else if (userStore.isLoggedIn && (to.path === '/login' || to.path === '/register')) {
    // 已登录用户访问登录/注册页，重定向到首页
    next('/teams')
  } else {
    next()
  }
})

export default router

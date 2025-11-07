
import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/lib/api'

export interface User {
  id: number
  username: string
  email: string
  nickname?: string
  avatar?: string
  role?: string
}

export const useUserStore = defineStore('user', () => {
  const token = ref<string | null>(localStorage.getItem('accessToken') || null)
  const currentUser = ref<User | null>(null)
  const isLoggedIn = ref(!!token.value)

  function setToken(t: string | null) {
    token.value = t
    if (t) localStorage.setItem('accessToken', t)
    else localStorage.removeItem('accessToken')
    isLoggedIn.value = !!t
  }
  function getToken() {
    return token.value
  }
  async function login(email: string, password: string) {
  try {
    const res = await api.post('/auth/login', { email, password })
    if (res.data?.token) { // 改成 token
      setToken(res.data.token)
      currentUser.value = res.data.user
      return true
    }
    return false
  } catch (err: any) {
    throw err?.response?.data?.error || '登录失败'
  }
}

  async function register(username: string, email: string, password: string, nickname: string) {
  try {
    const res = await api.post('/auth/register', { username, email, password, nickname })
    if (res.data?.user) {
      return res.data.user
    }
    return null
  } catch (err: any) {
    throw err?.response?.data?.error || '注册失败'
  }
}


  // 页面刷新时自动恢复登录状态
 if (token.value && !currentUser.value) {
  api.get('/user/profile').then(res => {
    currentUser.value = res.data
    isLoggedIn.value = true
  }).catch(() => {
    logout()
  })
}

function logout() {
  setToken(null)
  currentUser.value = null
}


  return { token, currentUser, isLoggedIn, setToken, getToken, login, register, logout }
})

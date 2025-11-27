
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
  phone?: string
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
  async function login(username: string, password: string) {
    try {
      const res = await api.post('/auth/login', { username, password })
      if (res.data?.accessToken) {
        setToken(res.data.accessToken)
        currentUser.value = res.data.user
        isLoggedIn.value = true
        return true
      }
      return false
    } catch (err: any) {
      const errorMsg = err?.response?.data?.error || err?.message || '登录失败'
      throw errorMsg
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
  async function restoreUser() {
    if (token.value && !currentUser.value) {
      try {
        const res = await api.get('/user/profile')
        if (res.data?.user) {
          currentUser.value = res.data.user
          isLoggedIn.value = true
        }
      } catch {
        logout()
      }
    }
  }

  function logout() {
    setToken(null)
    currentUser.value = null
    isLoggedIn.value = false
  }

  async function uploadAvatar(file: File) {
    const formData = new FormData()
    formData.append('avatar', file)
    try {
      const res = await api.post('/user/avatar', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      if (res.data?.url && currentUser.value) {
        // 后端返回相对路径，这里组合成完整的URL或者直接用相对路径（如果配置了baseURL）
        // 假设后端返回的是 /uploads/avatars/xxx.jpg
        // 我们需要确保它能被正确访问。如果 api baseURL 包含 host，这里可能需要调整。
        // 暂时直接使用返回的 url，通常后端返回的应该是可访问路径。
        // 如果是全路径最好，如果是相对路径，<img> 标签也能识别（相对于当前域名）。
        // 但是 api 请求是到后端端口 (e.g. 8080)，前端是 5173。
        // 所以如果返回相对路径，前端 img src="/uploads/..." 会请求前端服务器。
        // 前端 vite dev server 需要配置 proxy 或者后端返回全路径。
        // 更好的做法是后端返回完整 URL，或者前端拼接 API_BASE_URL。
        // 先简单处理，假设前端有 proxy 或者我们拼一下。
        
        const avatarUrl = res.data.url;
        currentUser.value.avatar = avatarUrl
        return avatarUrl
      }
    } catch (err: any) {
      throw err?.response?.data?.error || '头像上传失败'
    }
  }

  async function updateProfile(data: Partial<User>) {
    try {
      const res = await api.put('/user/profile', data)
      if (res.data?.user) {
        currentUser.value = res.data.user
        return res.data.user
      }
    } catch (err: any) {
      throw err?.response?.data?.error || '更新资料失败'
    }
  }

  // 初始化时恢复用户状态
  restoreUser()

  return { 
    token, 
    currentUser, 
    isLoggedIn, 
    setToken, 
    getToken, 
    login, 
    register, 
    logout, 
    restoreUser,
    uploadAvatar,
    updateProfile
  }
})

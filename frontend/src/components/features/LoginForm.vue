/**
 * 登录表单组件
 * 处理用户登录并在成功时存储token
 */
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'

const email = ref('')
const password = ref('')
const errorMsg = ref('')
const loading = ref(false)

const router = useRouter()
const userStore = useUserStore()

const handleSubmit = async (e: Event) => {
  e.preventDefault()
  if (loading.value) return
  loading.value = true
  errorMsg.value = ''
  try {
    const user = await userStore.login(email.value, password.value)
    if (user) {
      router.push('/')  // 登录成功
    } else {
      errorMsg.value = '登录失败，请检查账号和密码'
    }
  } catch (err: any) {
    errorMsg.value = typeof err === 'string' ? err : '登录过程中出现错误'
  } finally {
    loading.value = false
  }
}

</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <Card class="w-full max-w-md p-6">
      <h1 class="text-2xl font-bold mb-6">登录</h1>
      
      <form @submit="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">邮箱</label>
          <input
            v-model="email"
            type="email"
            required
            class="w-full p-2 border rounded-md"
            placeholder="请输入邮箱"
          />
        </div>
        
        <div>
          <label class="block text-sm font-medium mb-1">密码</label>
          <input
            v-model="password"
            type="password"
            required
            class="w-full p-2 border rounded-md"
            placeholder="请输入密码"
          />
        </div>

        <div v-if="errorMsg" class="text-red-500 text-sm">
          {{ errorMsg }}
        </div>

        <Button
          type="submit"
          :disabled="loading"
          class="w-full"
        >
          {{ loading ? '登录中...' : '登录' }}
        </Button>

        <div class="text-center mt-4">
          <router-link to="/register" class="text-blue-500 hover:text-blue-700">
            还没有账号？立即注册
          </router-link>
        </div>
      </form>
    </Card>
  </div>
</template>
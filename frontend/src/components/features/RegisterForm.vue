/**
 * 注册表单组件
 * 处理用户注册
 */
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const errorMsg = ref('')
const loading = ref(false)

const router = useRouter()
const userStore = useUserStore()

const handleSubmit = async (e: Event) => {
  e.preventDefault()
  if (loading.value) return
  if (password.value !== confirmPassword.value) {
    errorMsg.value = '两次输入的密码不一致'
    return
  }
  loading.value = true
  errorMsg.value = ''
  try {
    const user = await userStore.register(
      username.value,
      email.value,
      password.value,
      '' // 可选：nickname
    )
    if (user) {
      router.push('/login')  // 注册成功后跳转登录
    } else {
      errorMsg.value = '注册失败，请稍后重试'
    }
  } catch (err: any) {
    errorMsg.value = typeof err === 'string' ? err : '注册过程中出现错误'
  } finally {
    loading.value = false
  }
}

</script>

<template>
  <div class="min-h-screen flex items-center justify-center">
    <Card class="w-full max-w-md p-6">
      <h1 class="text-2xl font-bold mb-6">注册</h1>
      
      <form @submit="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">用户名</label>
          <input
            v-model="username"
            type="text"
            required
            minlength="3"
            maxlength="20"
            class="w-full p-2 border rounded-md"
            placeholder="请输入用户名"
          />
        </div>

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
            minlength="6"
            class="w-full p-2 border rounded-md"
            placeholder="请输入密码"
          />
        </div>

        <div>
          <label class="block text-sm font-medium mb-1">确认密码</label>
          <input
            v-model="confirmPassword"
            type="password"
            required
            minlength="6"
            class="w-full p-2 border rounded-md"
            placeholder="请再次输入密码"
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
          {{ loading ? '注册中...' : '注册' }}
        </Button>

        <div class="text-center mt-4">
          <router-link to="/login" class="text-blue-500 hover:text-blue-700">
            已有账号？立即登录
          </router-link>
        </div>
      </form>
    </Card>
  </div>
</template>
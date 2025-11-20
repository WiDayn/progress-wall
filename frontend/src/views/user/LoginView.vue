<template>
  <div class="min-h-screen flex items-center justify-center bg-background px-4">
    <Card class="w-full max-w-md">
      <CardHeader>
        <h1 class="text-2xl font-bold text-center">登录</h1>
        <p class="text-sm text-muted-foreground text-center mt-2">
          使用用户名或邮箱登录您的账户
        </p>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div class="space-y-2">
            <label for="username" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              用户名或邮箱
            </label>
            <Input
              id="username"
              v-model="username"
              type="text"
              required
              autocomplete="username"
              placeholder="请输入用户名或邮箱"
            />
          </div>

          <div class="space-y-2">
            <label for="password" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              密码
            </label>
            <Input
              id="password"
              v-model="password"
              type="password"
              required
              autocomplete="current-password"
              placeholder="请输入密码"
            />
          </div>

          <div v-if="errorMessage" class="text-sm text-destructive bg-destructive/10 p-3 rounded-md">
            {{ errorMessage }}
          </div>

          <Button
            type="submit"
            :disabled="loading"
            class="w-full"
          >
            {{ loading ? '登录中...' : '登录' }}
          </Button>

          <div class="text-center text-sm">
            <span class="text-muted-foreground">还没有账号？</span>
            <router-link to="/register" class="text-primary hover:underline font-medium">
              立即注册
            </router-link>
          </div>
        </form>
      </CardContent>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import CardHeader from '@/components/ui/CardHeader.vue'
import CardContent from '@/components/ui/CardContent.vue'
import Input from '@/components/ui/Input.vue'

const username = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = ref(false)

const router = useRouter()
const userStore = useUserStore()

const handleLogin = async () => {
  if (loading.value) return
  
  errorMessage.value = ''
  loading.value = true

  try {
    const success = await userStore.login(username.value, password.value)
    if (success) {
      router.push('/')
    } else {
      errorMessage.value = '登录失败，请检查用户名和密码'
    }
  } catch (err: any) {
    errorMessage.value = typeof err === 'string' ? err : '登录过程中出现错误'
  } finally {
    loading.value = false
  }
}
</script>

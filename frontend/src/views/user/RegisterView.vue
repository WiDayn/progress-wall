<template>
  <div class="min-h-screen flex items-center justify-center bg-background px-4">
    <Card class="w-full max-w-md">
      <CardHeader>
        <h1 class="text-2xl font-bold text-center">注册</h1>
        <p class="text-sm text-muted-foreground text-center mt-2">
          创建新账户以开始使用
        </p>
      </CardHeader>
      <CardContent>
        <form @submit.prevent="handleRegister" class="space-y-4">
          <div class="space-y-2">
            <label for="username" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              用户名 <span class="text-destructive">*</span>
            </label>
            <Input
              id="username"
              v-model="username"
              type="text"
              required
              autocomplete="username"
              placeholder="请输入用户名"
              minlength="3"
              maxlength="20"
            />
          </div>

          <div class="space-y-2">
            <label for="email" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              邮箱 <span class="text-destructive">*</span>
            </label>
            <Input
              id="email"
              v-model="email"
              type="email"
              required
              autocomplete="email"
              placeholder="请输入邮箱地址"
            />
          </div>

          <div class="space-y-2">
            <label for="nickname" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              昵称
            </label>
            <Input
              id="nickname"
              v-model="nickname"
              type="text"
              autocomplete="nickname"
              placeholder="请输入昵称（可选）"
            />
          </div>

          <div class="space-y-2">
            <label for="password" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              密码 <span class="text-destructive">*</span>
            </label>
            <Input
              id="password"
              v-model="password"
              type="password"
              required
              autocomplete="new-password"
              placeholder="请输入密码（至少6位）"
              minlength="6"
            />
          </div>

          <div class="space-y-2">
            <label for="confirmPassword" class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              确认密码 <span class="text-destructive">*</span>
            </label>
            <Input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              required
              autocomplete="new-password"
              placeholder="请再次输入密码"
              minlength="6"
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
            {{ loading ? '注册中...' : '注册' }}
          </Button>

          <div class="text-center text-sm">
            <span class="text-muted-foreground">已有账号？</span>
            <router-link to="/login" class="text-primary hover:underline font-medium">
              立即登录
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
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const nickname = ref('')
const errorMessage = ref('')
const loading = ref(false)

const router = useRouter()
const userStore = useUserStore()

const handleRegister = async () => {
  if (loading.value) return

  // 验证密码匹配
  if (password.value !== confirmPassword.value) {
    errorMessage.value = '两次输入的密码不一致'
    return
  }

  errorMessage.value = ''
  loading.value = true

  try {
    const user = await userStore.register(
      username.value,
      email.value,
      password.value,
      nickname.value || ''
    )
    if (user) {
      router.push('/login')
    } else {
      errorMessage.value = '注册失败，请稍后重试'
    }
  } catch (err: any) {
    errorMessage.value = typeof err === 'string' ? err : '注册过程中出现错误'
  } finally {
    loading.value = false
  }
}
</script>

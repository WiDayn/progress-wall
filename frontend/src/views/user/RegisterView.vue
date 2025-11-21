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
              placeholder="请输入用户名（3-20个字符）"
              minlength="3"
              maxlength="20"
              @blur="validateUsernameField"
            />
            <p v-if="usernameError" class="text-sm text-destructive">{{ usernameError }}</p>
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
              @blur="validateEmailField"
            />
            <p v-if="emailError" class="text-sm text-destructive">{{ emailError }}</p>
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
              placeholder="请输入昵称（可选，最多50个字符）"
              maxlength="50"
              @blur="validateNicknameField"
            />
            <p v-if="nicknameError" class="text-sm text-destructive">{{ nicknameError }}</p>
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
              placeholder="请输入密码（至少6位，包含字母和数字）"
              minlength="6"
              maxlength="128"
              @blur="validatePasswordField"
            />
            <p v-if="passwordError" class="text-sm text-destructive">{{ passwordError }}</p>
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
              @blur="validateConfirmPasswordField"
            />
            <p v-if="confirmPasswordError" class="text-sm text-destructive">{{ confirmPasswordError }}</p>
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
import {
  validateUsername,
  validateEmail,
  validatePassword,
  validateNickname,
  sanitizeInput
} from '@/utils/validation'

const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const nickname = ref('')
const errorMessage = ref('')
const loading = ref(false)

// 字段级错误提示
const usernameError = ref('')
const emailError = ref('')
const passwordError = ref('')
const confirmPasswordError = ref('')
const nicknameError = ref('')

const router = useRouter()
const userStore = useUserStore()

// 实时验证用户名
const validateUsernameField = () => {
  const result = validateUsername(username.value)
  usernameError.value = result.valid ? '' : (result.message || '')
  return result.valid
}

// 实时验证邮箱
const validateEmailField = () => {
  const result = validateEmail(email.value)
  emailError.value = result.valid ? '' : (result.message || '')
  return result.valid
}

// 实时验证密码
const validatePasswordField = () => {
  const result = validatePassword(password.value)
  passwordError.value = result.valid ? '' : (result.message || '')
  return result.valid
}

// 实时验证确认密码
const validateConfirmPasswordField = () => {
  if (!confirmPassword.value) {
    confirmPasswordError.value = ''
    return true
  }
  if (password.value !== confirmPassword.value) {
    confirmPasswordError.value = '两次输入的密码不一致'
    return false
  }
  confirmPasswordError.value = ''
  return true
}

// 实时验证昵称
const validateNicknameField = () => {
  const result = validateNickname(nickname.value)
  nicknameError.value = result.valid ? '' : (result.message || '')
  return result.valid
}

const handleRegister = async () => {
  if (loading.value) return

  // 清除所有错误
  errorMessage.value = ''
  usernameError.value = ''
  emailError.value = ''
  passwordError.value = ''
  confirmPasswordError.value = ''
  nicknameError.value = ''

  // 验证所有字段
  const isUsernameValid = validateUsernameField()
  const isEmailValid = validateEmailField()
  const isPasswordValid = validatePasswordField()
  const isConfirmPasswordValid = validateConfirmPasswordField()
  const isNicknameValid = validateNicknameField()

  if (!isUsernameValid || !isEmailValid || !isPasswordValid || !isConfirmPasswordValid || !isNicknameValid) {
    errorMessage.value = '请检查并修正表单中的错误'
    return
  }

  loading.value = true

  try {
    // 清理输入数据，防止XSS
    const sanitizedUsername = sanitizeInput(username.value)
    const sanitizedEmail = sanitizeInput(email.value)
    const sanitizedNickname = nickname.value ? sanitizeInput(nickname.value) : ''

    const user = await userStore.register(
      sanitizedUsername,
      sanitizedEmail,
      password.value, // 密码不需要HTML转义，但需要确保不包含特殊字符
      sanitizedNickname
    )
    if (user) {
      router.push('/login')
    } else {
      errorMessage.value = '注册失败，请稍后重试'
    }
  } catch (err: any) {
    const errorMsg = typeof err === 'string' ? err : err?.response?.data?.error || '注册过程中出现错误'
    errorMessage.value = errorMsg
  } finally {
    loading.value = false
  }
}
</script>

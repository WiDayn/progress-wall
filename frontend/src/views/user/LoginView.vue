<template>
  <div class="login-container">
    <h1>Login</h1>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label for="email">Email</label>
        <input v-model="email" type="email" id="email" placeholder="Enter your email" required />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input v-model="password" type="password" id="password" placeholder="Enter your password" required />
      </div>

      <button type="submit">Login</button>
    </form>

    <p class="error-message" v-if="errorMessage">{{ errorMessage }}</p>

    <p>
      Don't have an account?
      <a href="/register">Register here</a>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import http from '@/utils/http'

const email = ref('')
const password = ref('')
const errorMessage = ref('')

const handleLogin = async () => {
  errorMessage.value = ''
  try {
    const response = await http.post('/auth/login', {
      email: email.value,
      password: password.value
    })

    // 保存 JWT 到 localStorage
    localStorage.setItem('jwt', response.data.token)

    alert('Login successful!')

    // 登录成功后跳转主页（可修改为你的主页路径）
    window.location.href = '/dashboard'
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Login failed'
  }
}
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 100px auto;
  padding: 24px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.form-group {
  margin-bottom: 16px;
}

label {
  display: block;
  margin-bottom: 4px;
  font-weight: bold;
}

input {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

button {
  width: 100%;
  padding: 10px;
  font-weight: bold;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #66b1ff;
}

.error-message {
  color: red;
  margin-top: 12px;
}
</style>

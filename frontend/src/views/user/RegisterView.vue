<template>
  <div class="register-container">
    <h1>Register</h1>
    <form @submit.prevent="handleRegister">
      <div class="form-group">
        <label for="username">Username</label>
        <input v-model="username" type="text" id="username" placeholder="Enter your username" required />
      </div>

      <div class="form-group">
        <label for="email">Email</label>
        <input v-model="email" type="email" id="email" placeholder="Enter your email" required />
      </div>

      <div class="form-group">
        <label for="nickname">Nickname</label>
        <input v-model="nickname" type="text" id="nickname" placeholder="Enter your nickname" />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input v-model="password" type="password" id="password" placeholder="Enter your password" required />
      </div>

      <button type="submit">Register</button>
    </form>

    <p class="error-message" v-if="errorMessage">{{ errorMessage }}</p>

    <p>
      Already have an account?
      <a href="/login">Login here</a>
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import http from '@/utils/http'

const username = ref('')
const email = ref('')
const password = ref('')
const nickname = ref('')
const errorMessage = ref('')

const handleRegister = async () => {
  errorMessage.value = ''
  try {
    await http.post('/auth/register', {
      username: username.value,
      email: email.value,
      password: password.value,
      nickname: nickname.value
    })

    alert('Registration successful! You can now log in.')

    // 注册成功后跳转登录页
    window.location.href = '/login'
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Registration failed'
  }
}
</script>

<style scoped>
.register-container {
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
  background-color: #67c23a;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

button:hover {
  background-color: #85ce61;
}

.error-message {
  color: red;
  margin-top: 12px;
}
</style>

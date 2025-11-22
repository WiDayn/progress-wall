import axios from 'axios'

const http = axios.create({
  baseURL: 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json'
  }
})

// 添加请求拦截器
http.interceptors.request.use(config => {
  const token = localStorage.getItem('jwt')
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export default http

/**
 * 表单验证工具函数
 * 提供常用的输入验证规则，防止XSS攻击和无效输入
 */

/**
 * 转义HTML特殊字符，防止XSS攻击
 */
export function escapeHtml(text: string): string {
  const map: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  }
  return text.replace(/[&<>"']/g, (m) => map[m])
}

/**
 * 验证用户名格式
 * 要求：3-20个字符，只能包含字母、数字、下划线
 */
export function validateUsername(username: string): { valid: boolean; message?: string } {
  if (!username || username.trim().length === 0) {
    return { valid: false, message: '用户名不能为空' }
  }
  
  const trimmed = username.trim()
  if (trimmed.length < 3) {
    return { valid: false, message: '用户名长度至少3个字符' }
  }
  if (trimmed.length > 20) {
    return { valid: false, message: '用户名长度不能超过20个字符' }
  }
  
  const usernameRegex = /^[a-zA-Z0-9_]+$/
  if (!usernameRegex.test(trimmed)) {
    return { valid: false, message: '用户名只能包含字母、数字和下划线' }
  }
  
  return { valid: true }
}

/**
 * 验证邮箱格式
 */
export function validateEmail(email: string): { valid: boolean; message?: string } {
  if (!email || email.trim().length === 0) {
    return { valid: false, message: '邮箱不能为空' }
  }
  
  const emailRegex = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/
  if (!emailRegex.test(email.trim())) {
    return { valid: false, message: '请输入有效的邮箱地址' }
  }
  
  return { valid: true }
}

/**
 * 验证密码强度
 * 要求：至少6位，包含字母和数字，最大128位
 */
export function validatePassword(password: string): { valid: boolean; message?: string } {
  if (!password || password.length === 0) {
    return { valid: false, message: '密码不能为空' }
  }
  
  if (password.length < 6) {
    return { valid: false, message: '密码长度至少6位' }
  }
  
  if (password.length > 128) {
    return { valid: false, message: '密码长度不能超过128位' }
  }
  
  const hasLetter = /[a-zA-Z]/.test(password)
  const hasDigit = /[0-9]/.test(password)
  
  if (!hasLetter) {
    return { valid: false, message: '密码必须包含至少一个字母' }
  }
  
  if (!hasDigit) {
    return { valid: false, message: '密码必须包含至少一个数字' }
  }
  
  // 检查常见弱密码
  const weakPasswords = ['123456', 'password', '12345678', 'qwerty', 'abc123']
  if (weakPasswords.includes(password.toLowerCase())) {
    return { valid: false, message: '密码过于简单，请使用更复杂的密码' }
  }
  
  return { valid: true }
}

/**
 * 验证昵称格式
 * 要求：可选，如果填写则不超过50个字符
 */
export function validateNickname(nickname: string): { valid: boolean; message?: string } {
  if (!nickname || nickname.trim().length === 0) {
    return { valid: true } // 昵称是可选的
  }
  
  if (nickname.length > 50) {
    return { valid: false, message: '昵称长度不能超过50个字符' }
  }
  
  return { valid: true }
}

/**
 * 清理输入：去除首尾空格，转义HTML
 */
export function sanitizeInput(input: string): string {
  return escapeHtml(input.trim())
}


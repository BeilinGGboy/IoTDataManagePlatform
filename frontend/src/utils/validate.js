/**
 * 密码强度校验
 * @param {string} password
 * @returns {{ level: number, text: string, valid: boolean }}
 */
export function checkPasswordStrength(password) {
  if (!password) return { level: 0, text: '请输入密码', valid: false }

  let level = 0
  if (password.length >= 8) level++
  if (password.length >= 12) level++
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) level++
  if (/\d/.test(password)) level++
  if (/[!@#$%^&*(),.?":{}|<>]/.test(password)) level++

  const texts = ['弱', '较弱', '中等', '较强', '强', '很强']
  const valid = password.length >= 6

  return {
    level: Math.min(level, 5),
    text: texts[Math.min(level, 5)],
    valid,
  }
}

/**
 * 手机号校验
 */
export function isValidPhone(phone) {
  return /^1[3-9]\d{9}$/.test(phone)
}

/**
 * 邮箱校验
 */
export function isValidEmail(email) {
  if (!email) return false
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

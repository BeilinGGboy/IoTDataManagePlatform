import request from './request'

/**
 * 注册
 * @param {{ username: string, password: string, confirm_password: string, phone?: string, email?: string }} data
 */
export function register(data) {
  return request.post('/api/v1/auth/register', data)
}

/**
 * 登录
 * @param {{ username: string, password: string }} data
 */
export function login(data) {
  return request.post('/api/v1/auth/login', data)
}

/**
 * 获取当前用户信息（需鉴权）
 */
export function getMe() {
  return request.get('/api/v1/auth/me')
}

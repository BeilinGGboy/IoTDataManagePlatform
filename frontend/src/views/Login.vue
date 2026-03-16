<template>
  <div class="login-page">
    <div class="login-bg">
      <div class="login-bg-overlay" />
      <div class="login-bg-content">
        <h1>IoT 手表数据管理平台</h1>
        <p>智能设备数据采集与可视化分析</p>
      </div>
    </div>
    <div class="login-form-wrapper">
      <div class="login-form-card">
        <div class="login-header">
          <h2>欢迎登录</h2>
          <p>请选择登录方式</p>
        </div>

        <el-tabs v-model="activeTab" class="login-tabs">
          <!-- 密码登录 -->
          <el-tab-pane label="密码登录" name="password">
            <el-form
              ref="passwordFormRef"
              :model="passwordForm"
              :rules="passwordRules"
              @submit.prevent="handlePasswordLogin"
            >
              <el-form-item prop="username">
                <el-input
                  v-model="passwordForm.username"
                  placeholder="请输入用户名/手机号/邮箱"
                  size="large"
                  :prefix-icon="User"
                  clearable
                />
              </el-form-item>
              <el-form-item prop="password">
                <el-input
                  v-model="passwordForm.password"
                  type="password"
                  placeholder="请输入密码"
                  size="large"
                  :prefix-icon="Lock"
                  show-password
                  clearable
                  @keyup.enter="handlePasswordLogin"
                />
              </el-form-item>
              <el-form-item>
                <div class="form-options">
                  <el-checkbox v-model="passwordForm.remember">记住我</el-checkbox>
                  <el-link type="primary" :underline="false" @click="activeTab = 'forget'">
                    忘记密码？
                  </el-link>
                </div>
              </el-form-item>
              <el-form-item>
                <el-button
                  type="primary"
                  size="large"
                  :loading="loading"
                  class="login-btn"
                  @click="handlePasswordLogin"
                >
                  登 录
                </el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 短信验证码登录 -->
          <el-tab-pane label="短信登录" name="sms">
            <el-form
              ref="smsFormRef"
              :model="smsForm"
              :rules="smsRules"
              @submit.prevent="handleSmsLogin"
            >
              <el-form-item prop="phone">
                <el-input
                  v-model="smsForm.phone"
                  placeholder="请输入手机号"
                  size="large"
                  :prefix-icon="Iphone"
                  maxlength="11"
                  clearable
                />
              </el-form-item>
              <el-form-item prop="code">
                <div class="sms-code-input">
                  <el-input
                    v-model="smsForm.code"
                    placeholder="请输入验证码"
                    size="large"
                    :prefix-icon="Message"
                    maxlength="6"
                    clearable
                  />
                  <el-button
                    type="primary"
                    :disabled="countdown > 0 || !smsForm.phone"
                    @click="sendSmsCode"
                  >
                    {{ countdown > 0 ? `${countdown}s 后重发` : '获取验证码' }}
                  </el-button>
                </div>
              </el-form-item>
              <el-form-item>
                <el-button
                  type="primary"
                  size="large"
                  :loading="loading"
                  class="login-btn"
                  @click="handleSmsLogin"
                >
                  登 录
                </el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 忘记密码 -->
          <el-tab-pane label="忘记密码" name="forget">
            <el-form
              ref="forgetFormRef"
              :model="forgetForm"
              :rules="forgetRules"
              @submit.prevent="handleResetPassword"
            >
              <el-form-item prop="phone">
                <el-input
                  v-model="forgetForm.phone"
                  placeholder="请输入手机号"
                  size="large"
                  :prefix-icon="Iphone"
                  maxlength="11"
                  clearable
                />
              </el-form-item>
              <el-form-item prop="code">
                <div class="sms-code-input">
                  <el-input
                    v-model="forgetForm.code"
                    placeholder="请输入验证码"
                    size="large"
                    :prefix-icon="Message"
                    maxlength="6"
                    clearable
                  />
                  <el-button
                    type="primary"
                    :disabled="countdown > 0 || !forgetForm.phone"
                    @click="sendForgetCode"
                  >
                    {{ countdown > 0 ? `${countdown}s 后重发` : '获取验证码' }}
                  </el-button>
                </div>
              </el-form-item>
              <el-form-item prop="password">
                <el-input
                  v-model="forgetForm.password"
                  type="password"
                  placeholder="请输入新密码（6-20位）"
                  size="large"
                  :prefix-icon="Lock"
                  show-password
                  clearable
                />
                <div v-if="forgetForm.password" class="password-strength">
                  <el-progress
                    :percentage="strength.level * 20"
                    :stroke-width="6"
                    :color="strengthColor"
                  />
                  <span class="strength-text">密码强度：{{ strength.text }}</span>
                </div>
              </el-form-item>
              <el-form-item prop="confirmPassword">
                <el-input
                  v-model="forgetForm.confirmPassword"
                  type="password"
                  placeholder="请再次输入新密码"
                  size="large"
                  :prefix-icon="Lock"
                  show-password
                  clearable
                />
              </el-form-item>
              <el-form-item>
                <el-button
                  type="primary"
                  size="large"
                  :loading="loading"
                  class="login-btn"
                  @click="handleResetPassword"
                >
                  重置密码
                </el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>

        <div class="login-footer">
          <el-link type="info" :underline="false" @click="goBackToLogin">
            返回登录
          </el-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, Iphone, Message } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { checkPasswordStrength, isValidPhone } from '@/utils/validate'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const activeTab = ref('password')
const loading = ref(false)
const countdown = ref(0)
let countdownTimer = null

// 密码登录
const passwordFormRef = ref()
const passwordForm = reactive({
  username: '',
  password: '',
  remember: true,
})
const passwordRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

// 短信登录
const smsFormRef = ref()
const smsForm = reactive({
  phone: '',
  code: '',
})
const smsRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { validator: (_, v) => isValidPhone(v) || !v, message: '手机号格式不正确', trigger: 'blur' },
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' },
  ],
}

// 忘记密码
const forgetFormRef = ref()
const forgetForm = reactive({
  phone: '',
  code: '',
  password: '',
  confirmPassword: '',
})
const forgetRules = {
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { validator: (_, v) => isValidPhone(v) || !v, message: '手机号格式不正确', trigger: 'blur' },
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 6, message: '验证码为6位数字', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度6-20位', trigger: 'blur' },
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    {
      validator: (_, v) => v === forgetForm.password,
      message: '两次密码不一致',
      trigger: 'blur',
    },
  ],
}

// 密码强度
const strength = computed(() => checkPasswordStrength(forgetForm.password))
const strengthColor = computed(() => {
  const colors = ['#f56c6c', '#e6a23c', '#409eff', '#67c23a', '#67c23a']
  return colors[strength.value.level] || '#409eff'
})

// 切换 tab 时重置倒计时显示
watch(activeTab, () => {
  if (activeTab.value === 'password') {
    countdown.value = 0
  }
})

function goBackToLogin() {
  activeTab.value = 'password'
}

function startCountdown() {
  countdown.value = 60
  countdownTimer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) clearInterval(countdownTimer)
  }, 1000)
}

async function sendSmsCode() {
  if (!isValidPhone(smsForm.phone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }
  // 模拟发送：后端未实现，前端占位
  ElMessage.success('验证码已发送（演示模式：输入任意6位数字）')
  startCountdown()
}

async function sendForgetCode() {
  if (!isValidPhone(forgetForm.phone)) {
    ElMessage.warning('请输入正确的手机号')
    return
  }
  ElMessage.success('验证码已发送（演示模式：输入任意6位数字）')
  startCountdown()
}

async function handlePasswordLogin() {
  await passwordFormRef.value?.validate().catch(() => {})
  if (!passwordForm.username || !passwordForm.password) return

  loading.value = true
  try {
    // 后端未实现登录接口，前端模拟：任意账号密码均可进入
    userStore.setToken('demo-token-' + Date.now())
    userStore.setUserInfo({ username: passwordForm.username })
    ElMessage.success('登录成功')
    router.push(route.query.redirect || '/dashboard')
  } finally {
    loading.value = false
  }
}

async function handleSmsLogin() {
  await smsFormRef.value?.validate().catch(() => {})
  if (!smsForm.phone || !smsForm.code) return

  loading.value = true
  try {
    userStore.setToken('demo-token-sms-' + Date.now())
    userStore.setUserInfo({ phone: smsForm.phone })
    ElMessage.success('登录成功')
    router.push(route.query.redirect || '/dashboard')
  } finally {
    loading.value = false
  }
}

async function handleResetPassword() {
  await forgetFormRef.value?.validate().catch(() => {})
  if (!forgetForm.phone || !forgetForm.code || !forgetForm.password) return

  loading.value = true
  try {
    // 模拟重置成功
    ElMessage.success('密码重置成功，请使用新密码登录')
    activeTab.value = 'password'
    forgetForm.phone = ''
    forgetForm.code = ''
    forgetForm.password = ''
    forgetForm.confirmPassword = ''
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-page {
  display: flex;
  min-height: 100vh;
}

.login-bg {
  flex: 1;
  background: linear-gradient(135deg, #1e3a5f 0%, #0d2137 100%);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;

  .login-bg-overlay {
    position: absolute;
    inset: 0;
    background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.03'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
    opacity: 0.5;
  }

  .login-bg-content {
    position: relative;
    z-index: 1;
    text-align: center;
    color: #fff;

    h1 {
      font-size: 2rem;
      font-weight: 600;
      margin-bottom: 0.5rem;
    }
    p {
      font-size: 1rem;
      opacity: 0.8;
    }
  }
}

.login-form-wrapper {
  width: 440px;
  min-width: 440px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: #fff;
}

.login-form-card {
  width: 100%;
  max-width: 360px;
}

.login-header {
  margin-bottom: 1.5rem;
  h2 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1f2937;
  }
  p {
    font-size: 0.875rem;
    color: #6b7280;
    margin-top: 0.25rem;
  }
}

.login-tabs {
  :deep(.el-tabs__header) {
    margin-bottom: 1.5rem;
  }
  :deep(.el-tabs__item) {
    font-size: 0.9rem;
  }
  :deep(.el-tabs__nav-wrap::after) {
    display: none;
  }
}

.form-options {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sms-code-input {
  width: 100%;
  display: flex;
  gap: 12px;
  :deep(.el-input) {
    flex: 1;
  }
  .el-button {
    flex-shrink: 0;
    white-space: nowrap;
  }
}

.password-strength {
  margin-top: 8px;
  .strength-text {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
    display: block;
  }
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 1rem;
}

.login-footer {
  margin-top: 1.5rem;
  text-align: center;
}

@media (max-width: 900px) {
  .login-page {
    flex-direction: column;
  }
  .login-bg {
    min-height: 200px;
  }
  .login-form-wrapper {
    width: 100%;
    min-width: auto;
  }
}
</style>

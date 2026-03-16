<template>
  <div class="dashboard">
    <el-container>
      <el-aside width="220px" class="sidebar">
        <div class="logo">
          <span class="logo-icon">⌚</span>
          <span>IoT 数据平台</span>
        </div>
        <el-menu
          :default-active="route.path"
          router
          background-color="#1e3a5f"
          text-color="#a5b4c8"
          active-text-color="#fff"
        >
          <el-menu-item index="/dashboard">
            <el-icon><DataAnalysis /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
        </el-menu>
        <div class="sidebar-footer">
          <el-dropdown trigger="click" @command="handleUserCommand">
            <span class="user-info">
              <el-avatar :size="32">{{ userStore.userInfo?.username?.[0] || 'U' }}</el-avatar>
              <span class="username">{{ userStore.userInfo?.username || userStore.userInfo?.phone || '用户' }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-aside>
      <el-main class="main-content">
        <div class="page-header">
          <h1>仪表盘</h1>
          <p>欢迎使用 IoT 手表数据管理平台</p>
        </div>
        <el-card class="welcome-card">
          <p>前端已接入 Vue 3 + Element Plus，登录、短信验证、密码强度校验等功能已就绪。</p>
          <p>后端鉴权接口尚未实现，当前为演示模式，任意账号密码均可登录。</p>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { DataAnalysis } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

function handleUserCommand(cmd) {
  if (cmd === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  height: 100vh;
}

.el-container {
  height: 100%;
}

.sidebar {
  background: #1e3a5f;
  display: flex;
  flex-direction: column;

  .logo {
    height: 60px;
    display: flex;
    align-items: center;
    padding: 0 20px;
    color: #fff;
    font-weight: 600;
    font-size: 1rem;

    .logo-icon {
      font-size: 1.5rem;
      margin-right: 8px;
    }
  }

  .el-menu {
    flex: 1;
    border: none;
  }

  .sidebar-footer {
    padding: 16px;
    border-top: 1px solid rgba(255, 255, 255, 0.1);

    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      color: #a5b4c8;
      font-size: 0.875rem;

      .username {
        max-width: 120px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.main-content {
  background: #f5f7fa;
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;

  h1 {
    font-size: 1.5rem;
    font-weight: 600;
    color: #1f2937;
  }
  p {
    font-size: 0.875rem;
    color: #6b7280;
    margin-top: 4px;
  }
}

.welcome-card {
  max-width: 600px;

  p {
    margin-bottom: 8px;
    line-height: 1.6;
    color: #4b5563;
  }
  p:last-child {
    margin-bottom: 0;
  }
}
</style>

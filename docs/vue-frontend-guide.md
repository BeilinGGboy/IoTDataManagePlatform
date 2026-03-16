# Vue 前端框架说明

## 一、Vue 框架简介

Vue 是一个用于构建用户界面的渐进式 JavaScript 框架，具有以下特点：

| 概念 | 说明 |
|------|------|
| **响应式** | 数据变化时，界面自动更新，无需手动操作 DOM |
| **组件化** | 将页面拆成可复用的组件，每个组件包含模板、逻辑、样式 |
| **单文件组件 (.vue)** | 一个文件包含 `<template>`、`<script>`、`<style>`，结构清晰 |
| **组合式 API (Composition API)** | Vue 3 推荐，用 `setup()` 组织逻辑，便于复用和类型推导 |

### 与手写前端的对比

| 手写 HTML/JS | Vue |
|--------------|-----|
| 手动 `document.getElementById` 操作 DOM | 数据驱动，`v-model` 双向绑定 |
| 多处重复的 HTML 片段 | 组件复用 |
| 状态分散在各处 | Pinia 集中管理 |
| 路由用 hash 或手写 | Vue Router 声明式路由 |

---

## 二、项目结构

```
frontend/
├── index.html          # 入口 HTML
├── package.json        # 依赖
├── vite.config.js      # Vite 配置（构建、代理）
├── src/
│   ├── main.js         # 应用入口，挂载 Vue、插件
│   ├── App.vue         # 根组件
│   ├── router/         # 路由
│   │   └── index.js
│   ├── stores/         # Pinia 状态
│   │   └── user.js
│   ├── views/          # 页面组件
│   │   ├── Login.vue   # 登录页
│   │   └── Dashboard.vue
│   ├── api/            # 请求封装
│   │   └── request.js
│   └── utils/          # 工具函数
│       └── validate.js
└── public/             # 静态资源（直接复制到 dist）
    └── favicon.svg
```

---

## 三、技术栈

| 库 | 用途 |
|----|------|
| Vue 3 | 前端框架 |
| Vite | 构建工具，开发热更新快 |
| Vue Router | 路由 |
| Pinia | 状态管理（替代 Vuex） |
| Element Plus | UI 组件库 |
| Axios | HTTP 请求 |

---

## 四、前端功能说明

### 登录页 (Login.vue)

- **密码登录**：用户名/手机/邮箱 + 密码
- **短信登录**：手机号 + 6 位验证码（当前为演示，后端未实现）
- **忘记密码**：手机号 + 验证码 + 新密码，含密码强度校验

### 密码强度校验

- 长度、大小写、数字、特殊字符综合评分
- 弱 / 较弱 / 中等 / 较强 / 强

### 路由守卫

- 未登录访问需鉴权页面时，跳转到登录页
- 登录后跳回原目标页面

---

## 五、开发与构建

### 开发模式（前后端分离）

```bash
# 终端 1：启动后端
cd /Users/adai/Desktop/smartwatch-server
go run .

# 终端 2：启动前端
cd frontend
npm install
npm run dev
```

前端：http://localhost:5173  
后端：http://localhost:8080  

Vite 会将 `/api` 代理到后端，无需跨域。

### 生产构建

```bash
cd frontend
npm run build
```

产物在 `frontend/dist`，后端会优先从该目录提供静态资源。

---

## 六、与后端的对接

- 开发：前端 5173，后端 8080，Vite 代理 `/api` 到 8080
- 生产：前端构建后由 Gin 统一提供，同域访问，无跨域

后端鉴权接口尚未实现，当前为演示模式，任意账号密码均可登录。

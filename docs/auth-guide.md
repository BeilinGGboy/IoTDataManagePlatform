# 认证功能说明

## 一、功能概览

| 功能 | 状态 | 说明 |
|------|------|------|
| 注册 | ✅ | 用户名+密码，可选手机/邮箱 |
| 密码登录 | ✅ | 支持用户名/手机/邮箱登录 |
| 短信登录 | ⏳ | 暂未开放 |
| 忘记密码 | ⏳ | 暂未开放 |
| JWT 鉴权 | ✅ | Token 7 天有效 |
| 路由守卫 | ✅ | 未登录跳转登录页 |
| 401 处理 | ✅ | 自动清除 token 并跳转登录 |

## 二、前置条件

1. **数据库**：必须配置 MySQL/MariaDB，认证依赖 `admin_users` 表
2. **表结构**：执行 `scripts/init_db.sql` 或 `scripts/migrate_admin_users.sql`
3. **JWT 密钥**：生产环境在 `.env` 中设置 `JWT_SECRET`（建议 32 位以上随机字符串）

## 三、API 接口

### 注册

```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "admin",
  "password": "123456",
  "confirm_password": "123456",
  "phone": "13800138000",   // 选填
  "email": "admin@example.com"  // 选填
}

成功: 201, { "code": 0, "message": "注册成功", "data": { "token", "expires_at", "user" } }
失败: 400/409/500
```

### 登录

```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",  // 支持用户名/手机/邮箱
  "password": "123456"
}

成功: 200, { "code": 0, "message": "登录成功", "data": { "token", "expires_at", "user" } }
失败: 401
```

### 获取当前用户（需鉴权）

```
GET /api/v1/auth/me
Authorization: Bearer <token>

成功: 200, { "code": 0, "data": { "id", "username", "phone", "email", "role" } }
失败: 401
```

### 受保护接口

以下接口需在请求头携带 `Authorization: Bearer <token>`：

- `POST /api/v1/data/batch`
- `GET /api/v1/stats`

## 四、安全说明

- 密码使用 bcrypt 哈希存储
- JWT 默认 7 天过期
- 生产环境务必设置强随机 `JWT_SECRET`
- 用户名 3-32 位；密码 6-32 位；手机/邮箱格式校验

## 五、无数据库时的行为

未配置数据库时，服务以内存模式运行，**认证功能不可用**：

- `/api/v1/auth/*` 接口不注册
- `/api/v1/data/batch`、`/api/v1/stats` 不鉴权（兼容旧版）

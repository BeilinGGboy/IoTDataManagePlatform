# IoTDataManagePlatform / 智能手表数据服务器

智能手表数据采集与分析后端服务，基于 Gin 框架开发。这是第一个 Go 语言 IoT 和后端项目。

## 项目结构

```
smartwatch-server/
├── main.go                 # 主程序入口
├── go.mod                  # 依赖管理
├── go.sum                  # 依赖校验
│
├── api/                    # API 模块
│   ├── handlers/          # 请求处理器
│   │   └── data_handler.go
│   └── models/            # 数据模型
│       └── device_data.go
│
├── web/                    # 前端管理平台（可扩展）
│   ├── index.html         # 入口页
│   ├── css/style.css      # 样式
│   └── js/
│       ├── api.js         # API 调用
│       ├── router.js      # 路由（Hash 路由，支持多页）
│       └── app.js         # 页面渲染逻辑
│
├── config/                # 配置管理
│   └── database.go       # 数据库连接
└── README.md             # 项目说明
```

## 功能特性

- ✅ 接收批量设备数据上传
- ✅ 数据统计和监控
- ✅ 健康检查接口
- ✅ CORS 支持
- ✅ Web 管理平台（仪表盘、统计、多页可扩展）
- ✅ 数据库存储（MySQL）
- ⏳ 数据分析和查询（待实现）
- ⏳ 用户认证（待实现）

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

如遇网络问题，可设置国内代理：
```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 2. 运行服务器

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。**同一局域网内的其他设备**可通过 `http://<服务器IP>:8080` 访问管理平台。

### 3. 使用环境变量配置端口

```bash
PORT=3000 go run main.go
```

### 4. 手机/其他设备无法访问？

常见原因与排查：

| 问题 | 解决方案 |
|------|----------|
| **防火墙** | macOS：系统设置 → 网络 → 防火墙 → 允许 8080 端口；或临时关闭防火墙测试 |
| **不在同一网络** | 手机需连接与电脑相同的 WiFi（不能用流量） |
| **路由器隔离** | 部分路由器开启「AP 隔离」，会阻止设备互访，需在路由器设置中关闭 |

### 5. 域名访问（替代 IP）

同一局域网内，可用 **主机名.local** 访问（mDNS/Bonjour）：

- macOS：`http://你的电脑名.local:8080`（如 `http://MacBook-Pro.local:8080`）
- 需确保手机与电脑在同一 WiFi

### 6. 启用 HTTPS

默认使用 HTTP。如需 HTTPS（如部分功能要求安全上下文）：

```bash
# 生成自签名证书
chmod +x scripts/gen-cert.sh
./scripts/gen-cert.sh

# 使用证书启动
TLS_CERT=cert.pem TLS_KEY=key.pem go run main.go
```

自签名证书会导致浏览器提示「不安全」，需手动点击「高级」→「继续访问」。正式环境建议使用 Let's Encrypt 等可信证书。

### 7. 数据库配置（MySQL）

1. 安装 MySQL：`brew install mysql`，启动：`brew services start mysql`
2. 创建数据库：`mysql -u root -p < scripts/init_db.sql`
3. 配置连接：复制 `.env.example` 为 `.env`，填写 `DB_PASSWORD` 等
4. 详见 **[docs/MySQL安装与配置.md](docs/MySQL安装与配置.md)** 和 **[docs/数据库设计说明.md](docs/数据库设计说明.md)**

未配置数据库时，服务以内存模式运行（数据不持久化）。

### 8. 公网域名访问

从互联网通过域名访问，详见 **[docs/公网域名访问指南.md](docs/公网域名访问指南.md)**。

## API 接口

### 1. 批量上传数据

**接口:** `POST /api/v1/data/batch`

**请求体:**
```json
[
  {
    "device_id": "device-123456",
    "user_id": 1001,
    "timestamp": "2024-01-01T12:00:00Z",
    "heart_rate": {
      "bpm": 75,
      "timestamp": "2024-01-01T12:00:00Z"
    }
  },
  {
    "device_id": "device-123456",
    "user_id": 1001,
    "timestamp": "2024-01-01T12:01:00Z",
    "steps": {
      "steps": 1000,
      "distance": 700.0,
      "calories": 30,
      "timestamp": "2024-01-01T12:01:00Z"
    }
  }
]
```

**响应:**
```json
{
  "status": "success",
  "received": 2,
  "message": "Data received successfully"
}
```

### 2. 获取统计信息

**接口:** `GET /api/v1/stats`

**响应:**
```json
{
  "total_received": 100000,
  "total_batches": 1000,
  "duration": "1h30m0s",
  "avg_rate": 18.52,
  "uptime": "1h30m0s"
}
```

### 3. 健康检查

**接口:** `GET /health`

**响应:**
```json
{
  "status": "ok",
  "message": "Smartwatch data server is running"
}
```

## 部署到云服务器

### 阿里云 / 华为云部署步骤

1. **编译可执行文件**

```bash
# Linux 64位
GOOS=linux GOARCH=amd64 go build -o smartwatch-server main.go

# 或者直接在当前系统编译
go build -o smartwatch-server main.go
```

2. **上传到服务器**

```bash
scp smartwatch-server user@your-server-ip:/path/to/deploy/
```

3. **在服务器上运行**

```bash
# 设置端口（可选）
export PORT=8080

# 后台运行
nohup ./smartwatch-server > server.log 2>&1 &

# 或者使用 systemd（推荐）
```

4. **使用 systemd 管理（推荐）**

创建 `/etc/systemd/system/smartwatch-server.service`:

```ini
[Unit]
Description=Smartwatch Data Server
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/deploy
ExecStart=/path/to/deploy/smartwatch-server
Restart=always
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl start smartwatch-server
sudo systemctl enable smartwatch-server
```

5. **配置防火墙**

```bash
# 开放端口（以8080为例）
sudo ufw allow 8080/tcp
```

6. **使用 Nginx 反向代理（可选）**

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 开发计划

### 第一阶段（当前）
- [☑️] 基础 API 接口
- [☑️] 数据接收和统计
- [x] 健康检查

### 第二阶段
- [x] 数据库集成（MySQL）
- [x] 数据持久化存储
- [ ] 数据查询接口
- [ ] 用户认证（JWT）

### 第三阶段
- [ ] 数据分析和统计
- [ ] 实时数据推送（WebSocket）
- [ ] 数据可视化接口
- [ ] 异常检测和告警

### 第四阶段
- [ ] 微服务架构
- [ ] 分布式部署
- [ ] 监控和日志系统
- [ ] 性能优化

## 技术栈

- **Web 框架**: Gin
- **数据库**: PostgreSQL / SQLite（开发）
- **ORM**: GORM
- **部署**: Docker（待实现）

## 环境要求

- Go 1.21+
- 内存: 512MB+
- 磁盘: 1GB+（根据数据量）

## 注意事项

1. **生产环境配置**
   - 设置合适的端口
   - 配置 HTTPS（使用 Let's Encrypt）
   - 设置日志轮转
   - 配置监控告警

2. **性能优化**
   - 使用连接池
   - 批量插入数据库
   - 添加缓存层（Redis）
   - 负载均衡（多实例）

3. **安全**
   - 添加认证机制
   - 限制请求频率
   - 数据加密传输
   - SQL 注入防护

## 许可证

MIT License

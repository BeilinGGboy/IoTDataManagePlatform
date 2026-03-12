# 阿里云服务器部署指南

将 IoT 手表数据管理平台部署到阿里云，实现公网访问。

---

## 部署过程总结（2025-02-27 实测）

本次在阿里云轻量应用服务器（宝塔 Linux 面板镜像）上完成部署，主要步骤与注意事项如下：

| 步骤 | 说明 |
|------|------|
| 系统 | Aliyun Linux 3（RHEL 系），使用 `dnf`/`yum`，非 Ubuntu |
| 数据库 | 使用 **MariaDB**（MySQL 依赖冲突），root 为 unix_socket 认证，需创建 `iot_user` 供应用使用 |
| 防火墙 | 宝塔/轻量应用服务器：**应用搭建指南** → **站点设置** → **防火墙** → 添加 8080 端口 |
| 部署路径 | `/opt/iot-platform` |
| 访问地址 | `http://47.100.174.12:8080` |

---

## 一、前置准备

- 阿里云 ECS 或轻量应用服务器（Ubuntu 22.04 或 宝塔 Linux 面板 均可）
- 公网 IP
- 域名：1416153270107609.onaliyun.com（可选）
- 本地已安装 Go、MySQL/MariaDB 可运行项目

---

## 二、阿里云控制台配置

### 1. 安全组/防火墙放行端口（必做）

**方式 A：ECS 实例（有安全组）**

1. 登录 https://ecs.console.aliyun.com/
2. 左侧 **实例与镜像** → **实例** → 点击你的实例 ID
3. 在实例详情页找到 **安全组** → 点击安全组 ID
4. 点击 **入方向** → **手动添加**
5. 填写：协议类型 **TCP**，端口 **8080/8080**，授权对象 **0.0.0.0/0**
6. 保存

**方式 B：轻量应用服务器 / 宝塔镜像（无安全组入口）**

1. 控制台 → **应用搭建指南** → **站点设置**
2. 点击 **防火墙** → **去设置**
3. 添加规则：端口 **8080**，协议 **TCP**，来源 **0.0.0.0/0**

详见 **[docs/aliyun-security-group.md](aliyun-security-group.md)**

| 端口 | 协议 | 授权对象 | 说明 |
|------|------|----------|------|
| 22 | TCP | 0.0.0.0/0 | SSH |
| 8080 | TCP | 0.0.0.0/0 | 应用服务 |
| 80 | TCP | 0.0.0.0/0 | HTTP（可选，用于 Nginx） |

### 2. 域名解析（可选）

1. 阿里云控制台 → **域名** → 找到 `1416153270107609.onaliyun.com`
2. **解析设置** → **添加记录**：

| 记录类型 | 主机记录 | 记录值 | TTL |
|----------|----------|--------|-----|
| A | @ | 你的公网IP | 600 |

解析生效后，可通过 `http://1416153270107609.onaliyun.com:8080` 访问。

---

## 三、SSH 连接服务器

```bash
ssh root@你的公网IP
# 或使用密钥
ssh -i your-key.pem root@你的公网IP
```

首次连接会提示确认指纹，输入 `yes`。

---

## 四、服务器环境准备

### 1. 安装数据库（MySQL 或 MariaDB）

**Ubuntu：**

```bash
apt update
apt install -y mysql-server
systemctl start mysql
systemctl enable mysql
mysql_secure_installation
```

**Aliyun Linux / CentOS / RHEL（dnf/yum）：**

若 `mysql-server` 依赖冲突，可改用 MariaDB：

```bash
dnf install -y mariadb-server mariadb
systemctl start mariadb
systemctl enable mariadb
```

> MariaDB 的 root 默认使用 unix_socket 认证，无法用密码远程连接。应用应使用单独创建的 `iot_user`（见下文）。

### 2. 创建数据库和用户

```bash
# Ubuntu 用 mysql，Aliyun/CentOS 用 mariadb
mysql -u root -p   # 或 sudo mysql（MariaDB 无密码时）
```

在 MySQL/MariaDB 中执行：

```sql
CREATE DATABASE iot_watch_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 若使用 MariaDB，建议创建专用用户（避免 root 的 unix_socket 限制）
CREATE USER 'iot_user'@'localhost' IDENTIFIED BY '你的密码';
GRANT ALL PRIVILEGES ON iot_watch_db.* TO 'iot_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

然后导入表结构（需先上传 init_db.sql，见下文）：

```bash
# 使用 root 或 iot_user
mysql -u root -p iot_watch_db < /opt/iot-platform/scripts/init_db.sql
# 或
mysql -u iot_user -p iot_watch_db < /opt/iot-platform/scripts/init_db.sql
```

### 3. 安装 Go（用于在服务器编译，可选）

若在本地编译好上传，可跳过。

```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
source /etc/profile
```

---

## 五、上传代码并部署

### 方式 A：本地编译后上传（推荐）

**1. 在本地 Mac 编译**

```bash
cd /Users/adai/Desktop/smartwatch-server

# Linux 64 位
GOOS=linux GOARCH=amd64 go build -o smartwatch-server .
```

**2. 打包并上传**

```bash
# 使用脚本打包（推荐）
chmod +x scripts/deploy.sh
./scripts/deploy.sh

# 或手动打包（不含 .env，需在服务器上创建）
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")
GOOS=linux GOARCH=amd64 go build -ldflags="-X smartwatch-server/version.Version=$VERSION" -o smartwatch-server .
echo "$VERSION" > version.txt
tar -czvf deploy.tar.gz smartwatch-server web scripts/init_db.sql version.txt

# 上传（替换为你的公网 IP）
scp deploy.tar.gz root@你的公网IP:/opt/
```

**3. 在服务器解压并运行**

```bash
# SSH 登录后
mkdir -p /opt/iot-platform
cd /opt/iot-platform
tar -xzf /opt/deploy.tar.gz

# 创建 .env（需手动配置数据库密码）
cat > .env << 'EOF'
PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_USER=iot_user
DB_PASSWORD=你的数据库密码
DB_NAME=iot_watch_db
EOF
nano .env   # 修改 DB_PASSWORD（若用 MariaDB 的 iot_user，填 iot_user 的密码）

# 导入数据库表结构（若未创建库则先执行 mysql/mariadb 创建 iot_watch_db）
mysql -u iot_user -p iot_watch_db < scripts/init_db.sql

# 赋予执行权限
chmod +x smartwatch-server

# 测试运行
./smartwatch-server
```

### 方式 B：Git 拉取（若代码在 GitHub）

```bash
# 服务器上（Ubuntu 用 apt，Aliyun/CentOS 用 dnf）
apt install -y git   # 或 dnf install -y git
cd /opt
git clone https://github.com/BeilinGGboy/IoTDataManagePlatform.git iot-platform
cd iot-platform

# 安装 Go 后编译
go build -o smartwatch-server .

# 创建 .env（需手动配置）
nano .env
```

---

## 六、配置 .env

在服务器上编辑 `/opt/iot-platform/.env`：

```env
PORT=8080

DB_HOST=localhost
DB_PORT=3306
DB_USER=iot_user
DB_PASSWORD=你的数据库密码
DB_NAME=iot_watch_db
```

> 使用 MariaDB 时建议用 `iot_user`；使用 MySQL 且已配置 root 密码时可用 `root`。

---

## 七、使用 systemd 后台运行

创建服务文件：

```bash
nano /etc/systemd/system/smartwatch-server.service
```

内容：

```ini
[Unit]
Description=IoT Smartwatch Data Server
After=network.target mysql.service
# 若使用 MariaDB，改为：After=network.target mariadb.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/iot-platform
ExecStart=/opt/iot-platform/smartwatch-server
Restart=always
RestartSec=5
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
```

启动并开机自启：

```bash
systemctl daemon-reload
systemctl start smartwatch-server
systemctl enable smartwatch-server
systemctl status smartwatch-server
```

---

## 八、前台启动与后台启动

### 前台启动

在终端直接运行，日志输出到当前终端，关闭终端或按 `Ctrl+C` 会停止服务：

```bash
cd /opt/iot-platform
./smartwatch-server
```

适用于：调试、查看实时日志、临时测试。

### 后台启动（systemd 服务）

通过 systemd 管理，断开 SSH 后服务继续运行，开机自启：

```bash
systemctl start smartwatch-server
systemctl status smartwatch-server
```

适用于：生产环境、长期运行。

### 两者切换

| 操作 | 命令 |
|------|------|
| 后台 → 前台 | `systemctl stop smartwatch-server`，然后 `cd /opt/iot-platform && ./smartwatch-server` |
| 前台 → 后台 | 在前台运行时按 `Ctrl+C` 停止，再执行 `systemctl start smartwatch-server` |

> 同一时刻只能有一种方式在运行，切换前需先停止当前运行方式。

---

## 九、访问验证

- **IP 访问**：`http://47.100.174.12:8080`
- **域名访问**：`http://1416153270107609.onaliyun.com:8080`（需先完成域名解析）

---

## 十、版本管理

部署包中包含版本信息，便于确认服务器代码与本地一致。

### 查看版本

| 方式 | 命令/地址 |
|------|-----------|
| 接口查询 | `curl http://你的IP:8080/version` 或访问 `http://你的IP:8080/health`（含 version 字段） |
| 本地文件 | 解压后 `cat /opt/iot-platform/version.txt` |
| 启动日志 | 服务启动时在终端打印 `版本: xxx` |

### 版本号规则

- 有 tag：显示 tag（如 `v1.0.0`）
- 无 tag：显示短 commit hash（如 `abc1234`）
- 有未提交修改：后缀 `-dirty`

### 服务器拉取代码时记录版本

```bash
cd /opt/iot-platform
git pull
git rev-parse --short HEAD > .deployed_version   # 记录本次部署的 commit
# 编译时同样使用 deploy.sh 或手动加 -ldflags 注入版本
```

---

## 十二、常见问题

| 问题 | 处理 |
|------|------|
| 无法访问 | 检查安全组/防火墙是否放行 8080（宝塔镜像在站点设置→防火墙） |
| 数据库连接失败 | 检查 .env、MySQL/MariaDB 是否启动；MariaDB 建议用 iot_user 而非 root |
| MySQL 安装依赖冲突 | 改用 MariaDB：`dnf install -y mariadb-server mariadb` |
| 502/连接超时 | 确认服务已启动：`systemctl status smartwatch-server` |
| 域名无法访问 | 检查解析是否生效：`ping 1416153270107609.onaliyun.com` |

---

## 十二、去掉端口号（可选）

若希望用 `http://域名` 访问（不带 :8080），需配置 Nginx 反向代理：

```bash
apt install -y nginx
nano /etc/nginx/sites-available/default
```

在 `server` 块中添加：

```nginx
location / {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

然后 `systemctl restart nginx`，访问 `http://1416153270107609.onaliyun.com` 即可。

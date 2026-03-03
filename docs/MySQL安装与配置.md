# MySQL 安装与图形化工具配置指南

## 一、安装 MySQL（macOS）

### 1. 使用 Homebrew 安装

```bash
brew install mysql
```

### 2. 启动 MySQL 服务

```bash
brew services start mysql
```

### 3. 安全配置（设置 root 密码）

```bash
mysql_secure_installation
```

按提示操作：
- 设置 root 密码（建议为本项目单独设置，如 `iot123456`）
- 移除匿名用户：Y
- 禁止 root 远程登录：N（本地开发可保持）
- 删除 test 数据库：Y

### 4. 验证安装

```bash
mysql -u root -p
# 输入密码后进入 MySQL 命令行
```

---

## 二、图形化数据库管理工具

### 推荐一：DBeaver（免费、开源）

- 下载：https://dbeaver.io/download/
- 支持 MySQL、PostgreSQL、SQLite 等
- 功能完整：SQL 编辑、表设计、数据浏览、ER 图

### 推荐二：TablePlus（免费版可用）

- 下载：https://tableplus.com/
- 界面简洁，macOS 原生
- 免费版有连接数限制，个人开发足够

### 推荐三：MySQL Workbench（官方）

- 下载：https://dev.mysql.com/downloads/workbench/
- MySQL 官方工具，功能全面

---

## 三、创建项目数据库

### 方式一：命令行

```bash
mysql -u root -p
```

在 MySQL 中执行：

```sql
CREATE DATABASE iot_watch_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'iot_user'@'localhost' IDENTIFIED BY '你的密码';
GRANT ALL PRIVILEGES ON iot_watch_db.* TO 'iot_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

### 方式二：使用项目提供的脚本

在项目根目录执行（会提示输入 root 密码）：

```bash
cd /Users/adai/Desktop/smartwatch-server
mysql -u root -p < scripts/init_db.sql
```

若 `mysql` 命令找不到，可能是未加入 PATH。官方 DMG 安装的 MySQL 可尝试：

```bash
/usr/local/mysql/bin/mysql -u root -p < scripts/init_db.sql
```

---

## 四、在 DBeaver 中连接并查看数据库

### 1. 新建连接

1. 打开 DBeaver，点击左上角 **「新建连接」**（或 `Cmd + N`）
2. 选择 **MySQL**，点击「下一步」
3. 填写连接信息：

| 字段 | 值 |
|------|-----|
| 主机 | `localhost` |
| 端口 | `3306` |
| 数据库 | `iot_watch_db`（或留空先连接） |
| 用户名 | `root` |
| 密码 | 你的 MySQL root 密码 |

4. 点击 **「测试连接」**，若提示下载驱动则选择「下载」
5. 测试成功后点击「完成」

### 2. 查看数据库和表

1. 左侧 **数据库导航** 中展开连接
2. 展开 **「数据库」** → **「iot_watch_db」**
3. 展开 **「表」**，可看到：
   - `devices` - 设备表
   - `heart_rate_data` - 心率数据
   - `steps_data` - 步数数据
   - `sleep_data` - 睡眠数据
   - `sport_data` - 运动数据
   - `batch_uploads` - 批量上传记录

### 3. 查看表结构和数据

- **表结构**：右键表名 → 「查看表」或双击表名
- **数据**：右键表名 → 「查看数据」或选中表后按 `F4`
- **执行 SQL**：选中数据库 → 右键「SQL 编辑器」→「新建 SQL 脚本」

---

## 五、配置项目连接

在项目根目录创建 `.env` 文件（勿提交到 Git）：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=你的密码
DB_NAME=iot_watch_db
```

或使用环境变量启动：

```bash
DB_HOST=localhost DB_USER=root DB_PASSWORD=你的密码 DB_NAME=iot_watch_db go run main.go
```

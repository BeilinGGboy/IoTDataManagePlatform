#!/bin/bash
# 本地执行：打包项目并生成上传命令
# 用法: ./scripts/deploy.sh [服务器IP]

set -e
cd "$(dirname "$0")/.."

echo "=== 编译 Linux 版本 ==="
GOOS=linux GOARCH=amd64 go build -o smartwatch-server .

echo "=== 打包 ==="
tar -czvf deploy.tar.gz smartwatch-server web scripts/init_db.sql

echo ""
echo "=== 打包完成: deploy.tar.gz ==="
echo ""
echo "上传命令（替换 YOUR_SERVER_IP 为你的公网 IP）："
echo "  scp deploy.tar.gz root@YOUR_SERVER_IP:/opt/"
echo ""
echo "服务器上解压并运行："
echo "  mkdir -p /opt/iot-platform && cd /opt/iot-platform"
echo "  tar -xzf /opt/deploy.tar.gz"
echo "  # 创建 .env 并配置 DB_PASSWORD"
echo "  chmod +x smartwatch-server"
echo "  ./smartwatch-server"
echo ""
echo "详见: docs/aliyun-deployment-guide.md"

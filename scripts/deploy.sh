#!/bin/bash
# 本地执行：打包项目并生成上传命令
# 用法: ./scripts/deploy.sh [服务器IP]

set -e
cd "$(dirname "$0")/.."

# 版本号：优先用 git describe，否则用短 commit hash
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")
echo "=== 版本: $VERSION ==="

echo "=== 构建前端 (Vue) ==="
if [ -d frontend ]; then
  (cd frontend && npm ci --silent 2>/dev/null || npm install --silent) && (cd frontend && npm run build) || echo "⚠️ 前端构建跳过（无 frontend 或构建失败）"
else
  echo "⚠️ 无 frontend 目录，跳过前端构建"
fi

echo "=== 编译 Linux 版本 ==="
GOOS=linux GOARCH=amd64 go build -ldflags="-X smartwatch-server/version.Version=$VERSION" -o smartwatch-server .

echo "=== 打包 ==="
echo "$VERSION" > version.txt
if [ -d frontend/dist ]; then
  tar -czvf deploy.tar.gz smartwatch-server frontend/dist scripts/init_db.sql version.txt
else
  tar -czvf deploy.tar.gz smartwatch-server web scripts/init_db.sql version.txt
fi

echo ""
echo "=== 打包完成: deploy.tar.gz (版本 $VERSION) ==="
echo ""
echo "上传命令（替换为你的公网 IP）："
echo "  scp deploy.tar.gz root@47.100.174.12:/opt/"
echo ""
echo "服务器上解压并运行："
echo "  mkdir -p /opt/iot-platform && cd /opt/iot-platform"
echo "  tar -xzf /opt/deploy.tar.gz"
echo "  cat version.txt   # 查看本次部署版本"
echo "  # 创建 .env 并配置 DB_PASSWORD"
echo "  chmod +x smartwatch-server"
echo "  ./smartwatch-server"
echo ""
echo "详见: docs/aliyun-deployment-guide.md"

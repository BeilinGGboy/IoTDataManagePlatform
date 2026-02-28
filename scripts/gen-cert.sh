#!/bin/bash
# 生成自签名 HTTPS 证书（用于本地/局域网开发）
# 浏览器会提示「不安全」，需手动点击「高级」->「继续访问」

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OUT_DIR="${1:-$SCRIPT_DIR/../certs}"
mkdir -p "$OUT_DIR"

openssl req -x509 -newkey rsa:2048 -keyout "$OUT_DIR/key.pem" -out "$OUT_DIR/cert.pem" \
  -days 365 -nodes \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

echo "证书已生成: $OUT_DIR/cert.pem, $OUT_DIR/key.pem"
echo "启动 HTTPS: TLS_CERT=$OUT_DIR/cert.pem TLS_KEY=$OUT_DIR/key.pem go run main.go"

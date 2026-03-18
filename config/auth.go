package config

import "os"

// JWTSecret 返回 JWT 签名密钥，未设置时使用默认（生产环境务必设置 JWT_SECRET）
func JWTSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "iot-smartwatch-default-secret-change-in-production"
}

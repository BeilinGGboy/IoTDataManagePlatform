package auth

import (
	"smartwatch-server/api/models"
	"smartwatch-server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷，与 middleware 解析结构一致
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// IssueToken 签发 JWT
func IssueToken(u *models.AdminUser) (token string, expiresAt int64, err error) {
	secret := config.JWTSecret()
	expiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()
	claims := &Claims{
		UserID:   u.ID,
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString([]byte(secret))
	return
}

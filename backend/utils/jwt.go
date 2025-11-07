package utils

import (
	"time"

	"progress-wall-backend/config"

	"github.com/golang-jwt/jwt/v5"
)

// 生成 JWT token
func GenerateToken(userID uint, username string, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // 1天有效期
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))

}

package utils

import (
	"community-governance/application/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

const jwtSecretKey = "community-secret-key"

// GenerateJWT 生成 JWT Token
func GenerateJWT(userId string) (string, error) {
	claims := models.JWTClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "go-community",
		},
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 token
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	return signedToken, nil
}

// ParseJWT 解析并验证 JWT Token
func ParseJWT(tokenString string) (*models.JWTClaims, error) {
	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil // 使用硬编码的密钥进行解析
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	// 验证 token 是否有效
	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

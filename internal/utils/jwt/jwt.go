package jwt

import (
	"fmt"
	"time"

	"github.com/solumD/auth/internal/model"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken генерирует jwt-токен
func GenerateToken(info *model.UserInfo, secretKey []byte, duration time.Duration) (string, error) {
	if info == nil {
		return "", fmt.Errorf("info is nil")
	}
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		Username: info.Name,
		Role:     info.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

// VerifyToken валидирует jwt-токен и возвращает его claims
func VerifyToken(tokenStr string, secretKey []byte) (*model.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&model.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return secretKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %s", err.Error())
	}

	claims, ok := token.Claims.(*model.UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

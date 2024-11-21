package model

import "github.com/dgrijalva/jwt-go"

// UserInfo структура базовой информации о пользователе
type UserInfo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     int64  `json:"role"`
}

// UserClaims структура claims jwt-токена
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     int64  `json:"role"`
}

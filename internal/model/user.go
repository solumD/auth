package model

import (
	"database/sql"
	"time"
)

// User модель пользователя сервисного слоя
type User struct {
	ID              int64        `json:"id"`
	Name            string       `json:"name"`
	Email           string       `json:"email"`
	Password        string       `json:"password"`
	PasswordConfirm string       `json:"password_confirm"`
	Role            int64        `json:"role"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       sql.NullTime `json:"updated_at"`
}

// UserUpdate модель обновления пользователя сервисного слоя
type UserUpdate struct {
	ID    int64
	Name  *string
	Email *string
	Role  int64
}

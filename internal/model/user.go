package model

import (
	"database/sql"
	"time"
)

// Service user model
type User struct {
	ID              int64
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            int64
	CreatedAt       time.Time
	UpdatedAt       sql.NullTime
}

// Service user to update model
type UserUpdate struct {
	ID    int64
	Name  string
	Email string
	Role  int64
}

package model

import (
	"database/sql"
	"time"
)

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

type UserUpdate struct {
	ID    int64
	Name  string
	Email string
	Role  int64
}

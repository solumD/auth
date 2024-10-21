package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"username"`
	Email     string       `db:"email"`
	Password  string       `db:"password"`
	Role      int64        `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

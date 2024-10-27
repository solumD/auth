package model

// User модель пользователя в кэше
type User struct {
	ID          int64  `redis:"id"`
	Name        string `redis:"name"`
	Email       string `redis:"email"`
	Password    string `redis:"password"`
	Role        int64  `redis:"role"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at"`
}

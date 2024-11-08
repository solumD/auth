package cache

import (
	"context"

	"github.com/solumD/auth/internal/model"
)

// AuthCache интерфейс кэша
type AuthCache interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	DeleteUser(ctx context.Context, userID int64) error
}

package cache

import (
	"context"

	"github.com/solumD/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AuthCache интерфейс кэша
type AuthCache interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error)
}

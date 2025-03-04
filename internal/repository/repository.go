package repository

import (
	"context"

	"github.com/solumD/auth/internal/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

// UserRepository - интерфейс репо слоя user
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error)
}

// AuthRepository - интерфейс репо слоя auth
type AuthRepository interface {
	GetUser(ctx context.Context, name string) (*model.UserInfo, error)
}

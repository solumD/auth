package service

import (
	"context"

	"github.com/solumD/auth/internal/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

// UserService интерфейс сервисного слоя user
type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error)
}

// AuthService интерфейс сервисного слоя auth
type AuthService interface {
	Login(ctx context.Context, name string, password string) (string, error)
	GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

// AccessService интерфейс сервисного слоя access
type AccessService interface {
	Check(ctx context.Context, endpointAddress string) (*emptypb.Empty, error)
}

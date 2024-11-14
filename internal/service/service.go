package service

import (
	"context"

	"github.com/solumD/auth/internal/model"

	"google.golang.org/protobuf/types/known/emptypb"
)

// UserService - интерфейс сервисного слоя
type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	GetUser(ctx context.Context, userID int64) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error)
	DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error)
}

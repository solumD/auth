package user

import (
	"context"

	"github.com/solumD/auth/internal/api/user/errors"
	"github.com/solumD/auth/internal/converter"
	"github.com/solumD/auth/internal/logger"
	desc "github.com/solumD/auth/pkg/user_v1"

	"go.uber.org/zap"
)

// CreateUser - отправляет запрос в сервисный слой на создание пользователя
func (i *API) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	convertedReq := converter.ToUserFromDescUser(req)
	if convertedReq == nil {
		return nil, errors.ErrDescUserIsNil
	}

	userID, err := i.userService.CreateUser(ctx, convertedReq)
	if err != nil {
		return nil, err
	}

	logger.Info("inserted user", zap.Int64("userID", userID))

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

package user

import (
	"context"

	"github.com/solumD/auth/internal/api/user/errors"
	"github.com/solumD/auth/internal/converter"
	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/sys"
	"github.com/solumD/auth/internal/sys/codes"
	desc "github.com/solumD/auth/pkg/user_v1"

	"go.uber.org/zap"
)

// CreateUser - отправляет запрос в сервисный слой на создание пользователя
func (i *API) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	convertedReq := converter.ToUserFromDescUser(req)
	if convertedReq == nil {
		return nil, sys.NewCommonError(errors.ErrDescUserIsNil.Error(), codes.InvalidArgument)
	}

	userID, err := i.userService.CreateUser(ctx, convertedReq)
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.Aborted)
	}

	logger.Info("inserted user", zap.Int64("userID", userID))

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

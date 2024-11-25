package user

import (
	"context"

	"github.com/solumD/auth/internal/api/user/errors"
	"github.com/solumD/auth/internal/converter"
	"github.com/solumD/auth/internal/logger"
	desc "github.com/solumD/auth/pkg/user_v1"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser - отправляет запрос в сервисный слой на обновление данных пользователя
func (i *API) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	convertedReq := converter.ToUserFromDescUpdate(req)
	if convertedReq == nil {
		return nil, errors.ErrDescUserUpdateIsNil
	}

	_, err := i.userService.UpdateUser(ctx, convertedReq)
	if err != nil {
		return nil, err
	}

	logger.Info("updated user", zap.Int64("userID", req.GetId()))

	return &emptypb.Empty{}, nil
}

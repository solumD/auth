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

// GetUser - отправляет запрос в сервисный слой на получение данных пользователя
func (i *API) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	user, err := i.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.Aborted)
	}

	convertedUser := converter.ToDescUserFromService(user)
	if convertedUser == nil {
		return nil, sys.NewCommonError(errors.ErrUserModelIsNil.Error(), codes.Internal)
	}

	logger.Info("got user", zap.Int64("userID", req.GetId()))

	return convertedUser, nil
}

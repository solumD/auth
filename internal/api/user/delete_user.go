package user

import (
	"context"

	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/sys"
	"github.com/solumD/auth/internal/sys/codes"
	desc "github.com/solumD/auth/pkg/user_v1"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser - отправляет запрос в сервисный слой на удаление пользователя
func (i *API) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.Aborted)
	}

	logger.Info("deleted user", zap.Int64("userID", req.GetId()))

	return &emptypb.Empty{}, nil
}

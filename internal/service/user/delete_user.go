package user

import (
	"context"

	"github.com/solumD/auth/internal/logger"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser отправляет запрос в кэш, а затем в репо слой на удаление пользователя
func (s *srv) DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error) {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		_, errTx = s.userRepository.DeleteUser(ctx, userID)
		if errTx != nil {
			return errTx
		}

		errCache := s.authCache.DeleteUser(ctx, userID)
		if errCache != nil {
			logger.Error("failed to save user in cache", zap.Int64("userID", userID), zap.NamedError("error", errCache))
		} else {
			logger.Info("saved user in cache", zap.Int64("userID", userID))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

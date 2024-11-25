package user

import (
	"context"

	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/model"

	"go.uber.org/zap"
)

// GetUser отправляет запрос в кэш, а затем в репо слой на получение данных пользователя
func (s *srv) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.authCache.GetUser(ctx, userID)
	if err != nil {
		logger.Error("failed to get user from cache", zap.Int64("userID", userID))
	} else {
		logger.Info("got user from cache", zap.Int64("userID", userID))
		return user, nil
	}

	user, err = s.userRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = s.authCache.CreateUser(ctx, user)
	if err != nil {
		logger.Error("failed to save user in cache", zap.Int64("userID", userID), zap.NamedError("error", err))
	} else {
		logger.Info("saved user in cache", zap.Int64("userID", userID))
	}

	return user, nil
}

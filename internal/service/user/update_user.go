package user

import (
	"context"

	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/utils/validation"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser отправляет запрос в кэш, а затем в репо слой на обновление данных пользователя
func (s *srv) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	if user.Name != nil {
		err := validation.ValidateName(*user.Name)
		if err != nil {
			return nil, err
		}
	}

	if user.Email != nil {
		err := validation.ValidateEmail(*user.Email)
		if err != nil {
			return nil, err
		}
	}

	_, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	errCache := s.authCache.DeleteUser(ctx, user.ID)
	if errCache != nil {
		logger.Error("failed to delete user from cache", zap.Int64("userID", user.ID), zap.NamedError("error", errCache))
	} else {
		logger.Info("deleted user from cache", zap.Int64("userID", user.ID))
	}

	return &emptypb.Empty{}, nil
}

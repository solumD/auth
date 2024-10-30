package user

import (
	"context"

	"github.com/solumD/auth/internal/model"
)

// CreateUser отправляет запрос в репо слой на создание пользователя
func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	var userID int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		userID, errTx = s.authRepository.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}

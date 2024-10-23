package user

import (
	"context"

	"github.com/solumD/auth/internal/model"
)

func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	var userId int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		userId, errTx = s.authRepository.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return userId, nil
}

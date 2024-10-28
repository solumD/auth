package user

import (
	"context"
	"log"

	"github.com/solumD/auth/internal/model"
)

// CreateUser отправляет запрос в кэш, а затем в репо слой на создание пользователя
func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	var userID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		userID, errTx = s.authRepository.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		savedUser, errTx := s.authRepository.GetUser(ctx, userID)
		if errTx != nil {
			return errTx
		}

		errCache := s.authCache.CreateUser(ctx, savedUser)
		if errCache != nil {
			log.Printf("failed to save user %d in cache: %v ", userID, errCache)
		} else {
			log.Printf("saved user %d in cache", userID)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return userID, nil
}

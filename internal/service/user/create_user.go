package user

import (
	"context"
	"fmt"
	"log"

	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/validation"
)

// CreateUser отправляет запрос в репо слой на создание пользователя, а затем сохраняет данные в кэш
// CreateUser отправляет запрос в репо слой на создание пользователя, а затем сохраняет данные в кэш
func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	err := validation.ValidateName(user.Name)
	if err != nil {
		return 0, err
	}

	err = validation.ValidateEmail(user.Email)
	if err != nil {
		return 0, err
	}

	err = validation.ValidatePassword(user.Password)
	if err != nil {
		return 0, err
	}

	if user.Password != user.PasswordConfirm {
		return 0, fmt.Errorf("password and passwordConfirm do not match")
	}

	var userID int64

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
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

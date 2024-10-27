package user

import (
	"context"
	"log"

	"github.com/solumD/auth/internal/model"
)

// GetUser отправляет запрос в кэш, а затем в репо слой на получение данных пользователя
func (s *srv) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.authCache.GetUser(ctx, userID)
	if err != nil {
		log.Printf("failed to get user %d from cache", userID)
	} else {
		log.Printf("got user %d from cache", userID)
		return user, nil
	}

	user, err = s.authRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = s.authCache.CreateUser(ctx, user)
	if err != nil {
		log.Printf("failed to save user %d in cache", userID)
	} else {
		log.Printf("saved user %d in cache", userID)
	}

	return user, nil
}

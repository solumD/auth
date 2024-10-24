package user

import (
	"context"

	"github.com/solumD/auth/internal/model"
)

// GetUser - отправляет запрос в репо слой на получение данных пользователя
func (s *srv) GetUser(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.authRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

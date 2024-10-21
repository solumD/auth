package user

import (
	"context"

	"github.com/solumD/auth/internal/model"
)

func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	userId, err := s.authRepository.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

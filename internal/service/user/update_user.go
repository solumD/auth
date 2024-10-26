package user

import (
	"context"

	"github.com/solumD/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser отправляет запрос в репо слой на обновление данных пользователя
func (s *srv) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	_, err := s.authRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

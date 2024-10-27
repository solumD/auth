package user

import (
	"context"

	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/validation"

	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser отправляет запрос в кэш, а затем в репо слой на обновление данных пользователя
func (s *srv) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	err := validation.ValidateName(*user.Name)
	if err != nil {
		return nil, err
	}

	err = validation.ValidateEmail(*user.Email)
	if err != nil {
		return nil, err
	}

	_, err = s.authRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

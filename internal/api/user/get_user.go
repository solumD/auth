package user

import (
	"context"

	"github.com/solumD/auth/internal/api/user/errors"
	"github.com/solumD/auth/internal/converter"
	desc "github.com/solumD/auth/pkg/user_v1"
)

// GetUser - отправляет запрос в сервисный слой на получение данных пользователя
func (i *API) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	user, err := i.authService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	convertedUser := converter.ToDescUserFromService(user)
	if convertedUser == nil {
		return nil, errors.ErrUserModelIsNil
	}

	return convertedUser, nil
}

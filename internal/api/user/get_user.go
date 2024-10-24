package user

import (
	"context"

	"github.com/solumD/auth/internal/converter"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

// GetUser - отправляет запрос в сервисный слой на получение данных пользователя
func (i *Implementation) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	user, err := i.authService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return converter.ToDescUserFromService(user), nil
}

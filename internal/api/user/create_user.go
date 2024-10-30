package user

import (
	"context"
	"log"

	"github.com/solumD/auth/internal/converter"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

// CreateUser - отправляет запрос в сервисный слой на создание пользователя
func (i *AuthAPI) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	convertedReq, err := converter.ToUserFromDescUser(req)
	if err != nil {
		return nil, err
	}

	userID, err := i.authService.CreateUser(ctx, convertedReq)
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", userID)

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

package user

import (
	"context"
	"log"

	"github.com/solumD/auth/internal/api/user/errors"
	"github.com/solumD/auth/internal/converter"
	desc "github.com/solumD/auth/pkg/user_v1"
)

// CreateUser - отправляет запрос в сервисный слой на создание пользователя
func (i *API) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	convertedReq := converter.ToUserFromDescUser(req)
	if convertedReq == nil {
		return nil, errors.ErrDescUserIsNil
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

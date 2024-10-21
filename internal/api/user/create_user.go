package user

import (
	"context"
	"log"

	"github.com/solumD/auth/internal/converter"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	userID, err := i.authService.CreateUser(ctx, converter.ToUserFromDescUser(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", userID)

	return &desc.CreateUserResponse{
		Id: userID,
	}, nil
}

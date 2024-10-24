package user

import (
	"context"
	"log"

	desc "github.com/solumD/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser - отправляет запрос в сервисный слой на удаление пользователя
func (i *Implementation) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := i.authService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("deleted user with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

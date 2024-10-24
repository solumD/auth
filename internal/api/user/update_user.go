package user

import (
	"context"
	"log"

	"github.com/solumD/auth/internal/converter"
	desc "github.com/solumD/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser - отправляет запрос в сервисный слой на обновление данных пользователя
func (i *Implementation) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := i.authService.UpdateUser(ctx, converter.ToUserFromDescUpdate(req))
	if err != nil {
		return nil, err
	}

	log.Printf("updated user with id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

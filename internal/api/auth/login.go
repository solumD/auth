package auth

import (
	"context"
	"log"

	desc "github.com/solumD/auth/pkg/auth_v1"
)

// Login отправляет запрос в сервисный слой на авторизацию
func (a *API) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	token, err := a.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	log.Printf("succesful login for %s\n", req.GetUsername())

	return &desc.LoginResponse{
		RefreshToken: token,
	}, nil
}

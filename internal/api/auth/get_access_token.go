package auth

import (
	"context"
	"log"

	desc "github.com/solumD/auth/pkg/auth_v1"
)

// GetAccessToken отправляет запрос в сервисный слой на получение access токена
func (a *API) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	token, err := a.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	log.Println("got access token")

	return &desc.GetAccessTokenResponse{
		AccessToken: token,
	}, nil
}

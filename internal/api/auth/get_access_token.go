package auth

import (
	"context"

	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/sys"
	"github.com/solumD/auth/internal/sys/codes"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

// GetAccessToken отправляет запрос в сервисный слой на получение access токена
func (a *API) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	token, err := a.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.Unauthenticated)
	}

	logger.Info("got access token")

	return &desc.GetAccessTokenResponse{
		AccessToken: token,
	}, nil
}

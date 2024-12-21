package auth

import (
	"context"

	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/sys"
	"github.com/solumD/auth/internal/sys/codes"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

// GetRefreshToken отправляет запрос в сервисный слой на получение refresh токена
func (a *API) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	token, err := a.authService.GetRefreshToken(ctx, req.GetOldRefreshToken())
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.Unauthenticated)
	}

	logger.Info("got new refresh token")

	return &desc.GetRefreshTokenResponse{
		RefreshToken: token,
	}, nil
}

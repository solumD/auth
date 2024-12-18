package auth

import (
	"context"

	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/sys"
	"github.com/solumD/auth/internal/sys/codes"
	desc "github.com/solumD/auth/pkg/auth_v1"

	"go.uber.org/zap"
)

// Login отправляет запрос в сервисный слой на авторизацию
func (a *API) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := a.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.Unauthenticated)
	}

	accessToken, err := a.authService.GetAccessToken(ctx, refreshToken)
	if err != nil {
		return nil, sys.NewCommonError(err.Error(), codes.Unauthenticated)
	}

	logger.Info("succesful login", zap.String("username", req.GetUsername()))

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

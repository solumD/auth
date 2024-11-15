package auth

import (
	"context"
	"fmt"

	"github.com/solumD/auth/internal/utils/jwt"
)

// GetAccessToken принимает refresh token и на его основе создает и возвращает access token
func (s *srv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := jwt.VerifyToken(refreshToken, s.authCfg.RefreshTokenSecretKey())
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	userInfo, err := s.authRepository.GetUser(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(userInfo, s.authCfg.AccessTokenSecretKey(), s.authCfg.AccessTokenExp())
	if err != nil {
		return "", err
	}

	return token, nil
}

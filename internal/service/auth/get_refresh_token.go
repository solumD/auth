package auth

import (
	"context"
	"fmt"

	"github.com/solumD/auth/internal/utils/jwt"
)

// GetRefreshToken принимает старый refresh token и на его основе создает и возвращает новый
func (s *srv) GetRefreshToken(ctx context.Context, oldRefreshToken string) (string, error) {
	claims, err := jwt.VerifyToken(oldRefreshToken, s.authCfg.RefreshTokenSecretKey())
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	userInfo, err := s.authRepository.GetUser(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	token, err := jwt.GenerateToken(userInfo, s.authCfg.RefreshTokenSecretKey(), s.authCfg.RefreshTokenExp())
	if err != nil {
		return "", err
	}

	return token, nil
}

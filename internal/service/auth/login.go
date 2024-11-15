package auth

import (
	"context"
	"fmt"

	"github.com/solumD/auth/internal/utils/jwt"
	"github.com/solumD/auth/internal/utils/validation"
)

// Login валидирует данные пользователя, и если все ок, возвращает refresh token
func (s *srv) Login(ctx context.Context, username string, password string) (string, error) {
	err := validation.ValidateName(username)
	if err != nil {
		return "", err
	}

	err = validation.ValidatePassword(password)
	if err != nil {
		return "", err
	}

	userInfo, err := s.authRepository.GetUser(ctx, username)
	if err != nil {
		return "", err
	}

	if userInfo.Password != password {
		return "", fmt.Errorf("invalid password")
	}

	token, err := jwt.GenerateToken(userInfo, s.authCfg.RefreshTokenSecretKey(), s.authCfg.RefreshTokenExp())
	if err != nil {
		return "", err
	}

	return token, nil
}

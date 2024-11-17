package auth

import (
	"github.com/solumD/auth/internal/config"
	"github.com/solumD/auth/internal/repository"
	"github.com/solumD/auth/internal/service"
)

type srv struct {
	authRepository repository.AuthRepository
	authCfg        config.AuthConfig
}

// NewService возвращает новый объект сервисного слоя auth
func NewService(authRepo repository.AuthRepository, authCfg config.AuthConfig) service.AuthService {
	return &srv{
		authRepository: authRepo,
		authCfg:        authCfg,
	}
}

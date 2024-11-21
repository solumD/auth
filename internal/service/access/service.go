package access

import (
	"github.com/solumD/auth/internal/config"
	"github.com/solumD/auth/internal/service"
)

type srv struct {
	userAccesses map[string]struct{}
	authConfig   config.AuthConfig
}

// NewService возвращает новый объект сервисного слоя access
func NewService(userAccesses map[string]struct{}, authConfig config.AuthConfig) service.AccessService {
	return &srv{
		userAccesses: userAccesses,
		authConfig:   authConfig,
	}
}

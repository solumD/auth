package user

import (
	"github.com/solumD/auth/internal/repository"
	"github.com/solumD/auth/internal/service"
)

type srv struct {
	authRepository repository.AuthRepository
}

func NewService(authRepository repository.AuthRepository) service.AuthService {
	return &srv{
		authRepository: authRepository,
	}
}

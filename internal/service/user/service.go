package user

import (
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/repository"
	"github.com/solumD/auth/internal/service"
)

type srv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

func NewService(authRepository repository.AuthRepository, txManager db.TxManager) service.AuthService {
	return &srv{
		authRepository: authRepository,
		txManager:      txManager,
	}
}

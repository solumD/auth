package user

import (
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/repository"
	"github.com/solumD/auth/internal/service"
)

// Структура сервисного слоя с объектами репо слоя
// и транзакционного менеджера
type srv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

// NewService возвращает объект сервисного слоя
func NewService(authRepository repository.AuthRepository, txManager db.TxManager) service.AuthService {
	return &srv{
		authRepository: authRepository,
		txManager:      txManager,
	}
}

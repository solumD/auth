package user

import (
	"github.com/solumD/auth/internal/cache"
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/client/kafka"
	"github.com/solumD/auth/internal/repository"
	"github.com/solumD/auth/internal/service"
)

// Структура сервисного слоя с объектами репо слоя
// и транзакционного менеджера
type srv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
	authCache      cache.AuthCache
	kafkaProdcuer  kafka.Producer
}

// NewService возвращает объект сервисного слоя
func NewService(
	authRepository repository.AuthRepository,
	txManager db.TxManager,
	authCache cache.AuthCache,
	kafkaProducer kafka.Producer,
) service.AuthService {
	return &srv{
		authRepository: authRepository,
		txManager:      txManager,
		authCache:      authCache,
		kafkaProdcuer:  kafkaProducer,
	}
}

// NewMockService возвращает объект мока сервисного слоя
func NewMockService(deps ...interface{}) service.AuthService {
	serv := srv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.AuthRepository:
			serv.authRepository = s
		case cache.AuthCache:
			serv.authCache = s
		case db.TxManager:
			serv.txManager = s
		case kafka.Producer:
			serv.kafkaProdcuer = s
		}
	}

	return &serv
}

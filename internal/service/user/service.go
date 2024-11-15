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
	userRepository repository.UserRepository
	txManager      db.TxManager
	authCache      cache.AuthCache
	kafkaProducer  kafka.Producer
}

// NewService возвращает объект сервисного слоя
func NewService(
	userRepository repository.UserRepository,
	txManager db.TxManager,
	authCache cache.AuthCache,
	kafkaProducer kafka.Producer,
) service.UserService {
	return &srv{
		userRepository: userRepository,
		txManager:      txManager,
		authCache:      authCache,
		kafkaProducer:  kafkaProducer,
	}
}

// NewMockService возвращает объект мока сервисного слоя
func NewMockService(deps ...interface{}) service.UserService {
	serv := srv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.UserRepository:
			serv.userRepository = s
		case cache.AuthCache:
			serv.authCache = s
		case db.TxManager:
			serv.txManager = s
		case kafka.Producer:
			serv.kafkaProducer = s
		}
	}

	return &serv
}

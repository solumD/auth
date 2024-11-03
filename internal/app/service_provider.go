package app

import (
	"context"
	"log"

	authApi "github.com/solumD/auth/internal/api/user"
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/client/db/pg"
	"github.com/solumD/auth/internal/client/db/transaction"
	"github.com/solumD/auth/internal/closer"
	"github.com/solumD/auth/internal/config"
	"github.com/solumD/auth/internal/repository"
	authRepo "github.com/solumD/auth/internal/repository/user"
	"github.com/solumD/auth/internal/service"
	authSrv "github.com/solumD/auth/internal/service/user"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient  db.Client
	txManager db.TxManager

	authRepository repository.AuthRepository
	authService    service.AuthService
	authImpl       *authApi.AuthAPI
}

// NewServiceProvider возвращает новый объект API слоя
func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig инициализирует конфиг postgres
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GRPCConfig инициализирует конфиг grpc
func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// DBClient инициализирует клиент базы данных
func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create a db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// TxManager инициализирует менеджер транзакций
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// AuthRepository инициализирует репо слой
func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepo.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

// AuthService иницилизирует сервисный слой
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authSrv.NewService(s.AuthRepository(ctx), s.TxManager(ctx))
	}

	return s.authService
}

// AuthAPI инициализирует api слой
func (s *serviceProvider) AuthAPI(ctx context.Context) *authApi.AuthAPI {
	if s.authImpl == nil {
		s.authImpl = authApi.NewAuthAPI(s.AuthService(ctx))
	}

	return s.authImpl
}

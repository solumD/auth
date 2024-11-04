package app

import (
	"context"
	"log"

	redigo "github.com/gomodule/redigo/redis"
	authApi "github.com/solumD/auth/internal/api/user"
	authCache "github.com/solumD/auth/internal/cache"
	redisCache "github.com/solumD/auth/internal/cache/redis"
	"github.com/solumD/auth/internal/client/cache"
	"github.com/solumD/auth/internal/client/cache/redis"
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/client/db/pg"
	"github.com/solumD/auth/internal/client/db/transaction"
	"github.com/solumD/auth/internal/closer"
	"github.com/solumD/auth/internal/config"
	"github.com/solumD/auth/internal/repository"
	authRepoPG "github.com/solumD/auth/internal/repository/user"
	"github.com/solumD/auth/internal/service"
	authSrv "github.com/solumD/auth/internal/service/user"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	redisConfig config.RedisConfig

	dbClient    db.Client
	txManager   db.TxManager
	cacheClient cache.Client

	authCache      authCache.AuthCache
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

// RedisConfig инициализирует конфиг redis
func (s *serviceProvider) RedisConfig() config.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %v", err)
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
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
			log.Fatalf("postgres ping error: %v", err)
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

// RedisClient инициализирует клиент redis
func (s *serviceProvider) CacheClient(ctx context.Context) cache.Client {
	redisPool := &redigo.Pool{
		MaxIdle:     s.RedisConfig().MaxIdle(),
		IdleTimeout: s.RedisConfig().IdleTimeout(),
		DialContext: func(ctx context.Context) (redigo.Conn, error) {
			return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
		},
	}

	if s.cacheClient == nil {
		s.cacheClient = redis.NewClient(redisPool, s.RedisConfig())
	}

	err := s.cacheClient.Ping(ctx)
	if err != nil {
		log.Fatalf("redis ping error: %v", err)
	}

	return s.cacheClient
}

// AuthCache инициализирует кэш
func (s *serviceProvider) AuthCache(ctx context.Context) authCache.AuthCache {
	if s.authCache == nil {
		s.authCache = redisCache.NewRedisCache(s.CacheClient(ctx))
	}

	return s.authCache
}

// AuthRepository инициализирует репо слой
func (s *serviceProvider) AuthReposistory(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepoPG.NewRepository(s.DBClient(ctx))
		s.authRepository = authRepoPG.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

// AuthService иницилизирует сервисный слой
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authSrv.NewService(s.AuthReposistory(ctx), s.TxManager(ctx), s.AuthCache(ctx))
		s.authService = authSrv.NewService(s.AuthReposistory(ctx), s.TxManager(ctx), s.AuthCache(ctx))
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

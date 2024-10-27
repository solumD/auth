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

// Структура API слоя
type serviceProvider struct {
	pgConfig    config.PGConfig
	grpcConfig  config.GRPCConfig
	redisConfig config.RedisConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	authCache      authCache.AuthCache
	authRepository repository.AuthRepository
	authService    service.AuthService
	authImpl       *authApi.Implementation
}

// NewServiceProvider возвращает новый объект API слоя
func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

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

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient(ctx context.Context) cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	err := s.redisClient.Ping(ctx)
	if err != nil {
		log.Fatalf("redis ping error: %v", err)
	}

	return s.redisClient
}

func (s *serviceProvider) RedisCache(ctx context.Context) authCache.AuthCache {
	if s.authCache == nil {
		s.authCache = redisCache.NewRedisCache(s.RedisClient(ctx))
	}

	return s.authCache
}

func (s *serviceProvider) AuthReposistory(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepoPG.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authSrv.NewService(s.AuthReposistory(ctx), s.TxManager(ctx), s.RedisCache(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *authApi.Implementation {
	if s.authImpl == nil {
		s.authImpl = authApi.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}

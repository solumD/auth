package app

import (
	"context"
	"log"

	userApi "github.com/solumD/auth/internal/api/user"
	authCache "github.com/solumD/auth/internal/cache"
	redisCache "github.com/solumD/auth/internal/cache/redis"
	"github.com/solumD/auth/internal/client/cache"
	"github.com/solumD/auth/internal/client/cache/redis"
	"github.com/solumD/auth/internal/client/db"
	"github.com/solumD/auth/internal/client/db/pg"
	"github.com/solumD/auth/internal/client/db/transaction"
	"github.com/solumD/auth/internal/client/kafka"
	"github.com/solumD/auth/internal/client/kafka/producer"
	"github.com/solumD/auth/internal/closer"
	"github.com/solumD/auth/internal/config"
	"github.com/solumD/auth/internal/repository"
	userRepo "github.com/solumD/auth/internal/repository/user"
	"github.com/solumD/auth/internal/service"
	userSrv "github.com/solumD/auth/internal/service/user"

	redigo "github.com/gomodule/redigo/redis"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	redisConfig         config.RedisConfig
	httpConfig          config.HTTPConfig
	swaggerConfig       config.SwaggerConfig
	kafkaProducerConfig config.KafkaProducerConfig

	dbClient    db.Client
	txManager   db.TxManager
	cacheClient cache.Client

	authCache      authCache.AuthCache
	userRepository repository.UserRepository
	userService    service.UserService
	userImpl       *userApi.API

	kafkaProducer kafka.Producer
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

// HTTPConfig ининициализирует конфиг http сервера
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config")
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

// HTTPConfig ининициализирует конфиг http сервера
func (s *serviceProvider) SwaggerConfig() config.HTTPConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config")
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// KafkaProducerConfigининициализирует конфиг продюсера kafka
func (s *serviceProvider) KafkaProducerConfig() config.KafkaProducerConfig {
	if s.kafkaProducerConfig == nil {
		cfg, err := config.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka producer konfig")
		}

		s.kafkaProducerConfig = cfg
	}

	return s.kafkaProducerConfig
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

// KafkaProducer инициализрует продюсер kafka
func (s *serviceProvider) KafkaProducer(_ context.Context) kafka.Producer {
	if s.kafkaProducer == nil {
		p, err := producer.New(s.KafkaProducerConfig().Brokers(), s.KafkaProducerConfig().Config())
		if err != nil {
			log.Fatalf("failed to create kafka producer: %v", err)
		}

		closer.Add(p.Close)
		s.kafkaProducer = p
	}

	return s.kafkaProducer
}

// AuthCache инициализирует кэш
func (s *serviceProvider) AuthCache(ctx context.Context) authCache.AuthCache {
	if s.authCache == nil {
		s.authCache = redisCache.NewRedisCache(s.CacheClient(ctx))
	}

	return s.authCache
}

// UserRepository инициализирует репо слой
func (s *serviceProvider) UserReposistory(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// UserService иницилизирует сервисный слой
func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userSrv.NewService(s.UserReposistory(ctx), s.TxManager(ctx), s.AuthCache(ctx), s.KafkaProducer(ctx))
	}

	return s.userService
}

// UserAPI инициализирует api слой
func (s *serviceProvider) UserAPI(ctx context.Context) *userApi.API {
	if s.userImpl == nil {
		s.userImpl = userApi.NewAPI(s.UserService(ctx))
	}

	return s.userImpl
}

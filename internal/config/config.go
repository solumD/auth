package config

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

// GRPCConfig интерфейс grpc конфига
type GRPCConfig interface {
	Address() string
}

// LoggerConfig интерфейс конфига логгера
type LoggerConfig interface {
	Level() string
}

// PGConfig интерфейс postgres конфига
type PGConfig interface {
	DSN() string
}

// RedisConfig интерфейс redis конфига
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

// HTTPConfig интерфейс конфига http-сервера
type HTTPConfig interface {
	Address() string
}

// SwaggerConfig интерфейс конфига swagger http-сервера
type SwaggerConfig interface {
	Address() string
}

// PrometheusConfig интерфейс конфига prometheus http-сервера
type PrometheusConfig interface {
	Address() string
}

// KafkaProducerConfig интерфейс конфига продюсера kafka
type KafkaProducerConfig interface {
	Brokers() []string
	Config() *sarama.Config
}

// AuthConfig интерфейс конфига auth сервиса
type AuthConfig interface {
	RefreshTokenSecretKey() []byte
	AccessTokenSecretKey() []byte
	RefreshTokenExp() time.Duration
	AccessTokenExp() time.Duration
}

// AccessConfig интерфейс конфига access сервиса
type AccessConfig interface {
	UserAccessesMap() (map[string]struct{}, error)
}

// Load читает .env файл по указанному пути
// и загружает переменные в проект
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

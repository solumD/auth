package config

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

const (
	kafkaBrokersEnvName = "KAFKA_BROKERS"
)

type kafkaProducerConfig struct {
	brokers []string
}

// NewKafkaProducerConfig returns new kafka producer config
func NewKafkaProducerConfig() (KafkaProducerConfig, error) {
	brokersStr := os.Getenv(kafkaBrokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}

	brokers := strings.Split(brokersStr, ",")

	return &kafkaProducerConfig{
		brokers: brokers,
	}, nil
}

// Brokers returns list of broker's addresses
func (cfg *kafkaProducerConfig) Brokers() []string {
	return cfg.brokers
}

// Config returns sarama producer config
func (cfg *kafkaProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	return config
}
package producer

import (
	"github.com/solumD/auth/internal/client/kafka"

	"github.com/IBM/sarama"
)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

// New возвращает новый продюсер kafka
func New(brokersAddresses []string, cfg *sarama.Config) (kafka.Producer, error) {
	producer, err := sarama.NewSyncProducer(brokersAddresses, cfg)
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{
		producer: producer,
	}, nil
}

// SendMessage отправляет сообщение в брокер
func (p *kafkaProducer) SendMessage(msg *sarama.ProducerMessage) (int32, int64, error) {
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, err
	}

	return partition, offset, nil
}

// Close закрывает соединение продюсера
func (p *kafkaProducer) Close() error {
	if p.producer != nil {
		if err := p.producer.Close(); err != nil {
			return err
		}
	}

	return nil
}

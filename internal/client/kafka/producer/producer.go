package producer

import (
	"errors"

	"github.com/solumD/auth/internal/client/kafka"

	"github.com/IBM/sarama"
)

var (
	errProducerNotInitialized = errors.New("producer is not unitialized")
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
func (p *kafkaProducer) SendMessage(msg *sarama.ProducerMessage) *kafka.Response {
	if p.producer == nil {
		return &kafka.Response{
			Partition: 0,
			Offset:    0,
			Err:       errProducerNotInitialized,
		}
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return &kafka.Response{
			Partition: 0,
			Offset:    0,
			Err:       err,
		}
	}

	return &kafka.Response{
		Partition: partition,
		Offset:    offset,
		Err:       nil,
	}
}

// Close закрывает соединение продюсера
func (p *kafkaProducer) Close() error {
	if p.producer == nil {
		return errProducerNotInitialized
	}

	if err := p.producer.Close(); err != nil {
		return err
	}

	return nil
}

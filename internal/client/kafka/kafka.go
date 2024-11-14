package kafka

import "github.com/IBM/sarama"

// Producer интерфейс продюсера кафки
type Producer interface {
	SendMessage(msg *sarama.ProducerMessage) (int32, int64, error)
	Close() error
}

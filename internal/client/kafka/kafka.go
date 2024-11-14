package kafka

import "github.com/IBM/sarama"

// Producer интерфейс продюсера кафки
type Producer interface {
	SendMessage(msg *sarama.ProducerMessage) *Response
	Close() error
}

// Response ответ продюсера после вызова SendMessage
type Response struct {
	Partition int32
	Offset    int64
	Err       error
}

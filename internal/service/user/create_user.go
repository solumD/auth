package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/utils/validation"

	"github.com/IBM/sarama"
)

const (
	authTopicName = "auth"
)

// CreateUser отправляет запрос в репо слой на создание пользователя, а затем сохраняет данные в кэш
func (s *srv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	err := validation.ValidateName(user.Name)
	if err != nil {
		return 0, err
	}

	err = validation.ValidateEmail(user.Email)
	if err != nil {
		return 0, err
	}

	err = validation.ValidatePassword(user.Password)
	if err != nil {
		return 0, err
	}

	if user.Password != user.PasswordConfirm {
		return 0, fmt.Errorf("password and passwordConfirm do not match")
	}

	var userID int64

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		userID, errTx = s.userRepository.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		savedUser, errTx := s.userRepository.GetUser(ctx, userID)
		if errTx != nil {
			return errTx
		}

		errCache := s.authCache.CreateUser(ctx, savedUser)
		if errCache != nil {
			log.Printf("failed to save user %d in cache: %v\n", userID, errCache)
		} else {
			log.Printf("saved user %d in cache\n", userID)
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to marshall data: %v\n", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: authTopicName,
		Value: sarama.StringEncoder(data),
	}

	res := s.kafkaProducer.SendMessage(msg)
	if res.Err != nil {
		log.Printf("failed to send message in Kafka: %v\n", err)
	}

	log.Printf("message sent to partition %d with offset %d\n", res.Partition, res.Offset)

	return userID, nil
}

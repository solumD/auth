package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/solumD/auth/internal/logger"
	"github.com/solumD/auth/internal/model"
	"github.com/solumD/auth/internal/utils/validation"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
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
			logger.Error("failed to save user in cache", zap.Int64("userID", userID), zap.NamedError("error", errCache))
		} else {
			logger.Info("saved user in cache", zap.Int64("userID", userID))
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	data, err := json.Marshal(user)
	if err != nil {
		logger.Error("failed to marshall data", zap.NamedError("error", err))
	}

	msg := &sarama.ProducerMessage{
		Topic: authTopicName,
		Value: sarama.StringEncoder(data),
	}

	res := s.kafkaProducer.SendMessage(msg)
	if res.Err != nil {
		logger.Error("failed to send message in Kafka", zap.NamedError("error", err))
	}

	logger.Info("message sent in Kafka", zap.Int32("partition", res.Partition), zap.Int64("offset", res.Offset))

	return userID, nil
}

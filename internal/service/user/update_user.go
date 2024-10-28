package user

import (
	"context"
	"log"

	"github.com/solumD/auth/internal/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateUser отправляет запрос в кэш, а затем в репо слой на обновление данных пользователя
func (s *srv) UpdateUser(ctx context.Context, user *model.UserUpdate) (*emptypb.Empty, error) {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		_, errTx = s.authRepository.UpdateUser(ctx, user)

		if errTx != nil {
			return errTx
		}

		errCache := s.authCache.DeleteUser(ctx, user.ID)
		if errCache != nil {
			log.Printf("failed to delete user %d from cache: %v", user.ID, errCache)
		} else {
			log.Printf("deleted user %d from cache", user.ID)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

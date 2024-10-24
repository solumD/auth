package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser - отправляет запрос в репо слой на удаление пользователя
func (s *srv) DeleteUser(ctx context.Context, userID int64) (*emptypb.Empty, error) {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		_, errTx = s.authRepository.DeleteUser(ctx, userID)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

package converter

import (
	"fmt"

	"github.com/solumD/auth/internal/model"
	modelRepo "github.com/solumD/auth/internal/repository/user/model"
)

// ToUserFromRepo конвертирует модель пользователя репо слоя в
// модель сервисного слоя
func ToUserFromRepo(user *modelRepo.User) (*model.User, error) {
	if user == nil {
		return nil, fmt.Errorf("convertion failed, user model is nil")
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

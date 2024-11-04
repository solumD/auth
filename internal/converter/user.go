package converter

import (
	"github.com/solumD/auth/internal/model"
	desc "github.com/solumD/auth/pkg/auth_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToDescUserFromService конвертирует сервисную модель пользователя в
// в gRPC модель
func ToDescUserFromService(user *model.User) *desc.GetUserResponse {
	if user == nil {
		return nil
	}

	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// ToUserFromDescUser конвертирует модель пользователя API слоя в
// модель сервисного слоя
func ToUserFromDescUser(user *desc.CreateUserRequest) *model.User {
	if user == nil {
		return nil
	}

	return &model.User{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            int64(user.Role),
	}
}

// ToUserFromDescUpdate конвертирует модель обновления пользователя API слов в
// модель сервисного слоя
func ToUserFromDescUpdate(user *desc.UpdateUserRequest) *model.UserUpdate {
	if user == nil {
		return nil
	}

	u := &model.UserUpdate{}
	u.ID = user.Id
	u.Role = int64(user.Role)

	var name, email string

	if user.Name != nil {
		name = user.Name.GetValue()
		u.Name = &name
	}

	if user.Email != nil {
		email = user.Email.GetValue()
		u.Email = &email
	}

	return u
}

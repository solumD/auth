package converter

import (
	"github.com/solumD/auth/internal/model"
	desc "github.com/solumD/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user *model.User) *desc.GetUserResponse {
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

func ToUserFromDescUser(user *desc.CreateUserRequest) *model.User {
	return &model.User{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            int64(user.Role),
	}
}

func ToUserFromDescUpdate(user *desc.UpdateUserRequest) *model.UserUpdate {
	return &model.UserUpdate{
		ID:    user.GetId(),
		Name:  user.GetName().Value,
		Email: user.GetEmail().Value,
		Role:  int64(user.GetRole()),
	}
}
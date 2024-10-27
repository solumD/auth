package converter

import (
	"database/sql"
	"time"

	modelCache "github.com/solumD/auth/internal/cache/redis/model"
	"github.com/solumD/auth/internal/model"
)

// ToUserFromCache конвертирует модель пользователя из кэша
// в модель сервисного слоя
func ToUserFromCache(user *modelCache.User) *model.User {
	var updatedAt sql.NullTime
	if user.UpdatedAtNs != nil {
		updatedAt = sql.NullTime{
			Time:  time.Unix(0, *user.UpdatedAtNs),
			Valid: true,
		}
	}

	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		Password:  user.Password,
		CreatedAt: time.Unix(0, user.CreatedAtNs),
		UpdatedAt: updatedAt,
	}
}

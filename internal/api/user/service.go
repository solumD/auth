package user

import (
	"github.com/solumD/auth/internal/service"

	desc "github.com/solumD/auth/pkg/user_v1"
)

// API сруктура с заглушками gRPC-методов (при их отсутствии) и
// объект сервисного слоя (его интерфейса)
type API struct {
	desc.UnimplementedUserV1Server
	authService service.UserService
}

// NewAPI возвращает новый объект имплементации API-слоя
func NewAPI(authService service.UserService) *API {
	return &API{
		authService: authService,
	}
}

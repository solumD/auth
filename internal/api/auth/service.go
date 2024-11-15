package auth

import (
	"github.com/solumD/auth/internal/service"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

// API auth структура с заглушками gRPC-методов (при их отсутствии) и
// объект сервисного слоя (его интерфейса)
type API struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewAPI возвращает новый объект имплементации API-слоя auth
func NewAPI(authService service.AuthService) *API {
	return &API{
		authService: authService,
	}
}

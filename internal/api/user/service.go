package user

import (
	"github.com/solumD/auth/internal/service"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

// AuthAPI сруктура с заглушками gRPC-методов (при их отсутствии) и
// объект сервисного слоя (его интерфейса)
type AuthAPI struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewAuthAPI возвращает новый объект имплементации API-слоя
func NewAuthAPI(authService service.AuthService) *AuthAPI {
	return &AuthAPI{
		authService: authService,
	}
}

package user

import (
	"github.com/solumD/auth/internal/service"
	desc "github.com/solumD/auth/pkg/auth_v1"
)

// Implementation сруктура с заглушками gRPC-методов (при их отсутствии) и
// объект сервисного слоя (его интерфейса)
type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// NewImplementation возвращает новый объект имплементации API-слоя
func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}

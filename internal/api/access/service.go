package access

import (
	"github.com/solumD/auth/internal/service"

	desc "github.com/solumD/auth/pkg/access_v1"
)

// API access структура с заглушками gRPC-методов (при их отсутствии) и
// объект сервисного слоя (его интерфейса)
type API struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

// NewAPI возвращает новый объект API слоя access
func NewAPI(accessService service.AccessService) *API {
	return &API{
		accessService: accessService,
	}
}

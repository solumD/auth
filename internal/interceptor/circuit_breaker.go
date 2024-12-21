package interceptor

import (
	"context"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CircuitBreakerInterceptor контролирует допустимое rpc
type CircuitBreakerInterceptor struct {
	cb *gobreaker.CircuitBreaker
}

// NewCircuitBreakerInterceptor возвращает новый объект интерцептора метрик
func NewCircuitBreakerInterceptor(cb *gobreaker.CircuitBreaker) *CircuitBreakerInterceptor {
	return &CircuitBreakerInterceptor{
		cb: cb,
	}
}

// Unary ...
func (c *CircuitBreakerInterceptor) Unary(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := c.cb.Execute(func() (interface{}, error) {
		return handler(ctx, req)
	})

	if err != nil {
		if err == gobreaker.ErrOpenState {
			return nil, status.Error(codes.Unavailable, "service unavailable")
		}

		return nil, err
	}

	return res, nil
}

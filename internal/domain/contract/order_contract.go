package contract

import (
	"context"

	"github.com/order-nest/internal/domain"
)

// OrderRepository defines the data access interface for order operations.
type OrderRepository interface {
	Create(ctx context.Context, params domain.Order) (domain.CreateOrderResponse, error)
	List(ctx context.Context, parameters domain.OrderListFilter) (domain.OrderListResponse, error)
	Cancel(ctx context.Context, consignmentId string, userID uint64) error
}

// OrderUsecase defines the business logic interface for order operations.
type OrderUsecase interface {
	Create(ctx context.Context, params domain.CreateOrderRequest) (domain.CreateOrderResponse, error)
	List(ctx context.Context, parameters domain.OrderListFilter) (domain.OrderListResponse, error)
	Cancel(ctx context.Context, consignmentId string, userID uint64) error
}

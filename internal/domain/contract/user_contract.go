package contract

import (
	"context"

	"github.com/order-nest/internal/domain"
)

// UserRepository defines the data access interface for user operations.
type UserRepository interface {
	Create(context.Context, domain.CreateUserRequest) (domain.User, error)
	GetByUsername(context.Context, string) (domain.User, error)
	GetByID(context.Context, uint64) (domain.User, error)
}

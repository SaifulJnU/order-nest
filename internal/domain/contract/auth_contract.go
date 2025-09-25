package contract

import (
	"context"

	"github.com/order-nest/internal/domain"
)

// AuthUsecase defines the business logic interface for authentication-related operations.
type AuthUsecase interface {
	Login(ctx context.Context, request *domain.LoginRequest) (*domain.LoginResponse, error)
	CreateUser(ctx context.Context, user domain.CreateUserRequest) (domain.User, error)
}

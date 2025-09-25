package http

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/order-nest/config"
	"github.com/order-nest/database"
	"github.com/order-nest/internal/delivery/http/handler"
	"github.com/order-nest/internal/delivery/http/middleware"
	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/repository"
	"github.com/order-nest/internal/usecase"
	orderNestJwt "github.com/order-nest/pkg/auth_token"
	"github.com/order-nest/pkg/helper"
)

// Bootstrap initializes all dependencies and returns a ready-to-use Gin router
func Bootstrap(ctx context.Context) (*gin.Engine, error) {
	// Initialize validator
	helper.InitializeValidator()

	// Postgres DB connection
	db := database.ConnectPostgres()

	// JWT & Middleware
	jwtService := orderNestJwt.NewTokenService([]byte(config.GetConfig().JwtSecretKey))
	authMiddleware := middleware.NewAuth(jwtService)

	// Repositories & Usecases
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	// Migrations & default data
	if err := userRepo.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate users table: %w", err)
	}
	if err := orderRepo.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate orders table: %w", err)
	}

	// Test user
	authUsecase.CreateUser(ctx, domain.CreateUserRequest{
		Username: "01901901901@mailinator.com",
		Password: "321dsa",
	})

	// Router & Handlers (declarative route registry)
	r := gin.New()
	v1 := r.Group("/api/v1")

	authH := handler.NewAuthHandler(authUsecase)
	orderH := handler.NewOrderHandler(orderUsecase)

	// collect and register route tables
	register := func(routes interface {
		Routes(*middleware.Auth) []handler.RouteDef
	}) {
		for _, rd := range routes.Routes(authMiddleware) {
			switch rd.Method {
			case "GET":
				v1.GET(rd.Path, rd.Handlers...)
			case "POST":
				v1.POST(rd.Path, rd.Handlers...)
			case "PUT":
				v1.PUT(rd.Path, rd.Handlers...)
			case "DELETE":
				v1.DELETE(rd.Path, rd.Handlers...)
			case "PATCH":
				v1.PATCH(rd.Path, rd.Handlers...)
			case "HEAD":
				v1.HEAD(rd.Path, rd.Handlers...)
			case "OPTIONS":
				v1.OPTIONS(rd.Path, rd.Handlers...)
			default:
				// fallback to POST if unknown
				v1.Any(rd.Path, rd.Handlers...)
			}
		}
	}

	register(authH)
	register(orderH)

	return r, nil
}

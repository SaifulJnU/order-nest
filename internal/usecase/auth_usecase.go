package usecase

import (
	"context"
	"strconv"

	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/domain/contract"
	orderNestJwt "github.com/order-nest/pkg/auth_token"
	appLogger "github.com/order-nest/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo     contract.UserRepository
	tokenService *orderNestJwt.TokenService
}

func NewAuthUsecase(userRepo contract.UserRepository, tokenService *orderNestJwt.TokenService) contract.AuthUsecase {
	return &authUsecase{userRepo: userRepo, tokenService: tokenService}
}

func (a *authUsecase) CreateUser(ctx context.Context, user domain.CreateUserRequest) (domain.User, error) {
	appLogger.L().WithField("username", user.Username).Info("creating user")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		appLogger.L().WithError(err).Error("password hash failed")
		return domain.User{}, err
	}

	user.Password = string(hashedPassword)
	createdUser, err := a.userRepo.Create(ctx, user)
	if err != nil {
		appLogger.L().WithError(err).Error("user repository create failed")
		return domain.User{}, err
	}

	appLogger.L().WithField("user_id", createdUser.ID).Info("user created")
	return createdUser, nil
}

func (a *authUsecase) Login(ctx context.Context, request *domain.LoginRequest) (*domain.LoginResponse, error) {
	appLogger.L().WithField("username", request.Username).Info("login attempt")
	// Basic validation: username and password required
	if request.Username == "" || request.Password == "" {
		appLogger.L().Warn("login validation failed: empty username or password")
		return nil, domain.BadRequestError
	}

	user, err := a.userRepo.GetByUsername(ctx, request.Username)
	if err != nil {
		appLogger.L().WithError(err).Warn("get user by username failed")
		return nil, err
	}

	// Verify password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		appLogger.L().WithField("user_id", user.ID).Warn("password mismatch")
		return nil, domain.BadRequestError
	}

	// Generate tokens
	payload := orderNestJwt.Payload{Aud: strconv.FormatUint(user.ID, 10), Name: user.Username}
	token, err := a.tokenService.Generate(ctx, payload)
	if err != nil {
		appLogger.L().WithError(err).Error("token generation failed")
		return nil, domain.InternalServerError
	}

	appLogger.L().WithField("user_id", user.ID).Info("login success")
	return &domain.LoginResponse{
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

package usecase

import (
	"context"
	"strconv"

	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/domain/contract"
	orderNestJwt "github.com/order-nest/pkg/auth_token"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	user.Password = string(hashedPassword)
	createdUser, err := a.userRepo.Create(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}

func (a *authUsecase) Login(ctx context.Context, request *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Basic validation: username and password required
	if request.Username == "" || request.Password == "" {
		return nil, domain.BadRequestError
	}

	user, err := a.userRepo.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	// Verify password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, domain.BadRequestError
	}

	// Generate tokens
	payload := orderNestJwt.Payload{Aud: strconv.FormatUint(user.ID, 10), Name: user.Username}
	token, err := a.tokenService.Generate(ctx, payload)
	if err != nil {
		return nil, domain.InternalServerError
	}

	return &domain.LoginResponse{
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

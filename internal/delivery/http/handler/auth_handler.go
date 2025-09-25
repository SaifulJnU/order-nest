package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/order-nest/internal/domain/contract"
	appLogger "github.com/order-nest/pkg/logger"

	"github.com/gin-gonic/gin"
	customError "github.com/order-nest/internal/delivery/http/custom_error"
	"github.com/order-nest/internal/delivery/http/middleware"
	"github.com/order-nest/internal/domain"
)

// AuthHandler handles authentication-related HTTP endpoints.
type AuthHandler struct {
	authUsecase contract.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authUsecase contract.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Routes returns the route table for auth endpoints.
func (a *AuthHandler) Routes(authMiddleware *middleware.Auth) []RouteDef {
	return []RouteDef{
		{Method: http.MethodPost, Path: "/login", Handlers: []gin.HandlerFunc{a.login}},
		{Method: http.MethodGet, Path: "/logout", Handlers: []gin.HandlerFunc{authMiddleware.AuthRequired(a.logout)}},
	}
}

// login handles user authentication and token generation.
func (a *AuthHandler) login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customError.HTTPError{
			Message: "invalid payload",
			Type:    "error",
			Code:    http.StatusBadRequest,
		})
		return
	}

	authResponse, err := a.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		// Handle known domain errors
		if errors.Is(err, domain.BadRequestError) || errors.Is(err, domain.NotFoundError) {
			c.JSON(http.StatusBadRequest, customError.HTTPError{
				Message: "invalid credentials",
				Type:    "error",
				Code:    http.StatusBadRequest,
			})
			return
		}

		// Unknown internal error
		appLogger.L().WithError(err).Error("login failed")
		c.JSON(http.StatusInternalServerError, customError.HTTPError{
			Message: "internal server error",
			Type:    "error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, authResponse)
}

// Currently simulates delay and returns a basic response.
func (a *AuthHandler) logout(c *gin.Context) {
	time.Sleep(2 * time.Second)
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
		"type":    "success",
		"code":    http.StatusOK,
	})
}

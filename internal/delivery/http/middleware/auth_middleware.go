package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	customErr "github.com/order-nest/internal/delivery/http/custom_error"
	orderNestJwt "github.com/order-nest/pkg/auth_token"
	appLogger "github.com/order-nest/pkg/logger"
)

// Auth provides middleware for JWT authentication.
type Auth struct {
	tokenService *orderNestJwt.TokenService
}

// NewAuthMiddleware creates a new instance of Auth middleware.
func NewAuthMiddleware(tokenService *orderNestJwt.TokenService) *Auth {
	return &Auth{tokenService: tokenService}
}

// lookupBearerToken extracts a bearer token (case-insensitive) from Authorization header.
func lookupBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("authorization header must be Bearer <token>")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}

// AuthRequired is a Gin middleware that validates JWT and sets user info in the context.
func (a *Auth) AuthRequired(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token
		token, err := lookupBearerToken(c.Request)
		if err != nil {
			appLogger.L().WithError(err).Warn("auth: missing or invalid authorization header")
			c.JSON(http.StatusUnauthorized, customErr.Unauthrized)
			c.Abort()
			return
		}

		// Parse and validate token
		tokenPayload, err := a.tokenService.Parse(c.Request.Context(), token)
		if err != nil {
			appLogger.L().WithError(err).Warn("auth: token parse failed")
			c.JSON(http.StatusUnauthorized, customErr.Unauthrized)
			c.Abort()
			return
		}

		// Set token info in Gin context for downstream handlers
		c.Set("aud", tokenPayload.Aud)
		c.Set("username", tokenPayload.Name)

		next(c)
	}
}

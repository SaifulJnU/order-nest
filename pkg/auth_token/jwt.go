package auth_token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/order-nest/config"
)

// TokenService handles JWT generation and parsing
type TokenService struct {
	key []byte
}

// Payload represents the data stored in JWT claims
type Payload struct {
	Aud  string `json:"aud"`
	Name string `json:"name"`
}

// Token contains the generated access and refresh tokens
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Jti          string `json:"jti"`
}

// NewTokenService creates a new JWT service with the given key
func NewTokenService(key []byte) *TokenService {
	return &TokenService{key: key}
}

// Generate creates a new access and refresh token for the given payload
func (t *TokenService) Generate(ctx context.Context, payload Payload) (Token, error) {
	now := time.Now()
	jti := uuid.New().String()

	buildClaims := func(tokenType string, expSeconds time.Duration) jwt.MapClaims {
		return jwt.MapClaims{
			"aud":        payload.Aud,
			"name":       payload.Name,
			"jti":        jti,
			"iat":        now.Unix(),
			"exp":        now.Add(time.Second * expSeconds).Unix(),
			"token_type": tokenType,
		}
	}

	accessClaims := buildClaims("access", config.GetConfig().AccessTokenDuration)
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString(t.key)
	if err != nil {
		return Token{}, err
	}

	refreshClaims := buildClaims("refresh", config.GetConfig().RefreshTokenDuration)
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString(t.key)
	if err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(config.GetConfig().RefreshTokenDuration),
		Jti:          jti,
	}, nil
}

// Parse validates a token and extracts the payload
func (t *TokenService) Parse(ctx context.Context, tokenString string) (Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// Enforce HS256
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return t.key, nil
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, keyFunc, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return Payload{}, err
	}
	if !parsedToken.Valid {
		return Payload{}, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return Payload{}, errors.New("invalid token claims")
	}

	var p Payload
	if aud, ok := claims["aud"].(string); ok {
		p.Aud = aud
	}
	if name, ok := claims["name"].(string); ok {
		p.Name = name
	}
	return p, nil
}

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/order-nest/internal/delivery/http/middleware"
	"github.com/order-nest/internal/domain"
)

type mockAuthUsecase struct{ mock.Mock }

func (m *mockAuthUsecase) Login(_ context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	args := m.Called(req)
	if v := args.Get(0); v != nil {
		return v.(*domain.LoginResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockAuthUsecase) CreateUser(_ context.Context, _ domain.CreateUserRequest) (domain.User, error) {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.(domain.User), args.Error(1)
	}
	return domain.User{}, args.Error(1)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	m := new(mockAuthUsecase)
	h := NewAuthHandler(m)
	mw := &middleware.Auth{} // not used for login

	// route wiring as in Routes
	for _, rd := range h.Routes(mw) {
		if rd.Path == "/login" && rd.Method == http.MethodPost {
			r.POST(rd.Path, rd.Handlers...)
		}
	}

	reqBody := domain.LoginRequest{Username: "u", Password: "p"}
	token := &domain.LoginResponse{TokenType: "Bearer", AccessToken: "abc"}
	m.On("Login", &reqBody).Return(token, nil)

	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthHandler_Login_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	m := new(mockAuthUsecase)
	h := NewAuthHandler(m)
	mw := &middleware.Auth{}

	for _, rd := range h.Routes(mw) {
		if rd.Path == "/login" && rd.Method == http.MethodPost {
			r.POST(rd.Path, rd.Handlers...)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("{")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

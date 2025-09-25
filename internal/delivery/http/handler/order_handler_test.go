package handler

import (
	"bytes"
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

type mockOrderUsecase struct{ mock.Mock }

func (m *mockOrderUsecase) Create(_ context.Context, req domain.CreateOrderRequest) (domain.CreateOrderResponse, error) {
	args := m.Called(req)
	if v := args.Get(0); v != nil {
		return v.(domain.CreateOrderResponse), args.Error(1)
	}
	return domain.CreateOrderResponse{}, args.Error(1)
}

func (m *mockOrderUsecase) Cancel(_ context.Context, id string, uid uint64) error {
	args := m.Called(id, uid)
	return args.Error(0)
}

func (m *mockOrderUsecase) List(_ context.Context, f domain.OrderListFilter) (domain.OrderListResponse, error) {
	args := m.Called(f)
	if v := args.Get(0); v != nil {
		return v.(domain.OrderListResponse), args.Error(1)
	}
	return domain.OrderListResponse{}, args.Error(1)
}

func TestOrderHandler_Create_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	m := new(mockOrderUsecase)
	h := NewOrderHandler(m)
	_ = &middleware.Auth{}

	// Register route bypassing auth middleware; inject aud and call handler directly
	r.POST("/orders", func(c *gin.Context) {
		c.Set("aud", "1")
		h.createOrder(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte("{")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestOrderHandler_List_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	m := new(mockOrderUsecase)
	h := NewOrderHandler(m)
	_ = &middleware.Auth{}

	// Register route bypassing auth middleware; inject aud and call handler directly
	r.GET("/orders/all", func(c *gin.Context) {
		c.Set("aud", "1")
		h.listOrders(c)
	})

	// expected list call
	m.On("List", mock.AnythingOfType("domain.OrderListFilter")).Return(domain.OrderListResponse{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/all", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

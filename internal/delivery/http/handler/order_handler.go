package handler

import (
	"errors"
	"net/http"

	"github.com/order-nest/internal/domain/contract"
	appLogger "github.com/order-nest/pkg/logger"

	"github.com/gin-gonic/gin"
	orderNestError "github.com/order-nest/internal/delivery/http/custom_error"
	"github.com/order-nest/internal/delivery/http/middleware"
	"github.com/order-nest/internal/domain"
	"github.com/order-nest/pkg/helper"
)

// OrderHandler handles HTTP endpoints related to orders.
type OrderHandler struct {
	orderUsecase contract.OrderUsecase
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(orderUsecase contract.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

// Routes returns the route table for order endpoints.
func (o *OrderHandler) Routes(authMiddleware *middleware.Auth) []RouteDef {
	return []RouteDef{
		{Method: http.MethodPost, Path: "/orders", Handlers: []gin.HandlerFunc{authMiddleware.AuthRequired(o.createOrder)}},
		{Method: http.MethodPut, Path: "/orders/:consignment_id/cancel", Handlers: []gin.HandlerFunc{authMiddleware.AuthRequired(o.cancelOrder)}},
		{Method: http.MethodGet, Path: "/orders/all", Handlers: []gin.HandlerFunc{authMiddleware.AuthRequired(o.listOrders)}},
	}
}

// createOrder handles creating a new order.
func (o *OrderHandler) createOrder(c *gin.Context) {
	aud, ok := helper.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, orderNestError.Unauthrized)
		return
	}

	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, orderNestError.HTTPError{
			Message: "invalid payload",
			Type:    "error",
			Code:    http.StatusBadRequest,
		})
		return
	}
	req.CreatedBy = aud

	ord, err := o.orderUsecase.Create(c.Request.Context(), req)
	if err != nil {
		var validationErr *domain.ValidationError
		if errors.As(err, &validationErr) {
			c.JSON(http.StatusUnprocessableEntity, orderNestError.HTTPError{
				Message: "validation failed",
				Type:    "error",
				Code:    http.StatusUnprocessableEntity,
				Errors:  validationErr.ErrorMap,
			})
			return
		}
		appLogger.L().WithError(err).Error("order create failed")
		c.JSON(http.StatusInternalServerError, orderNestError.HTTPError{
			Message: "internal server error",
			Type:    "error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order Created Successfully",
		"type":    "success",
		"code":    http.StatusOK,
		"data":    ord,
	})
}

// cancelOrder handles cancelling an existing order by consignment ID.
func (o *OrderHandler) cancelOrder(c *gin.Context) {
	aud, ok := helper.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, orderNestError.Unauthrized)
		return
	}

	consignmentID := c.Param("consignment_id")
	if consignmentID == "" {
		c.JSON(http.StatusBadRequest, orderNestError.HTTPError{
			Message: "Bad request: missing consignment ID",
			Type:    "error",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := o.orderUsecase.Cancel(c.Request.Context(), consignmentID, aud); err != nil {
		c.JSON(http.StatusBadRequest, orderNestError.HTTPError{
			Message: "Unable to cancel order, please contact support",
			Type:    "error",
			Code:    http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order Cancelled Successfully",
		"type":    "success",
		"code":    http.StatusOK,
	})
}

// listOrders returns paginated orders created by the authenticated user.
func (o *OrderHandler) listOrders(c *gin.Context) {
	aud, ok := helper.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, orderNestError.Unauthrized)
		return
	}

	limit, page := helper.ParsePaginationParams(c)
	transferStatus := helper.GetUint8QueryParam(c, "transfer_status")
	archive := helper.GetUint8QueryParam(c, "archive")

	orders, err := o.orderUsecase.List(c.Request.Context(), domain.OrderListFilter{
		Limit:          int64(limit),
		Page:           int64(page),
		TransferStatus: transferStatus,
		Archive:        archive,
		CreatedBy:      aud,
	})
	if err != nil {
		appLogger.L().WithError(err).Error("order list failed")
		c.JSON(http.StatusInternalServerError, orderNestError.HTTPError{
			Message: "internal server error",
			Type:    "error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Orders successfully fetched.",
		"type":    "success",
		"code":    http.StatusOK,
		"data":    orders,
	})
}

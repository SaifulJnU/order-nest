package repository

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/repository/schema"
	appLogger "github.com/order-nest/pkg/logger"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository initializes the repository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create persists a new order.
func (o *OrderRepository) Create(ctx context.Context, params domain.Order) (domain.CreateOrderResponse, error) {
	appLogger.L().WithField("merchant_order_id", params.MerchantOrderId).Info("repo: creating order")
	repoOrder := mapOrderDomainToSchema(params)
	repoOrder.OrderConsignmentId = uuid.New().String()
	repoOrder.OrderCreatedAt = time.Now()
	repoOrder.UpdatedAt = time.Now()

	tx := o.db.WithContext(ctx)
	if err := tx.Create(&repoOrder).Error; err != nil {
		appLogger.L().WithError(err).Error("repo: order create failed")
		return domain.CreateOrderResponse{}, err
	}

	resp := domain.CreateOrderResponse{
		ConsignmentId:   repoOrder.OrderConsignmentId,
		MerchantOrderId: repoOrder.MerchantOrderId,
		OrderStatus:     repoOrder.OrderStatus,
		DeliveryFee:     repoOrder.DeliveryFee,
	}
	appLogger.L().WithField("consignment_id", resp.ConsignmentId).Info("repo: order created")
	return resp, nil
}

// Cancel sets an order status to canceled.
func (o *OrderRepository) Cancel(ctx context.Context, consignmentId string, userID uint64) error {
	appLogger.L().WithFields(map[string]interface{}{
		"consignment_id": consignmentId,
		"user_id":        userID,
	}).Info("repo: cancelling order")
	tx := o.db.WithContext(ctx)
	result := tx.Model(&schema.Order{}).
		Where("order_consignment_id = ?", consignmentId).
		Updates(map[string]interface{}{
			"order_status": domain.OrderStatusCanceled,
			"updated_by":   userID,
			"updated_at":   time.Now(),
		})
	if result.Error != nil {
		appLogger.L().WithError(result.Error).Error("repo: order cancel failed")
		return result.Error
	}
	appLogger.L().WithField("consignment_id", consignmentId).Info("repo: order cancelled")
	return nil
}

// List returns paginated orders for a user and filters.
func (o *OrderRepository) List(ctx context.Context, params domain.OrderListFilter) (domain.OrderListResponse, error) {
	appLogger.L().WithFields(map[string]interface{}{
		"created_by":      params.CreatedBy,
		"transfer_status": params.TransferStatus,
		"archive":         params.Archive,
		"limit":           params.Limit,
		"page":            params.Page,
	}).Info("repo: listing orders")
	limit, page, offset := calculatePagination(params.Limit, params.Page)

	tx := o.db.WithContext(ctx)
	base := tx.Model(&schema.Order{}).
		Where("created_by = ? AND transfer_status = ? AND archive = ?",
			params.CreatedBy, params.TransferStatus, params.Archive,
		)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		appLogger.L().WithError(err).Error("repo: count orders failed")
		return domain.OrderListResponse{}, err
	}

	var orders []schema.Order
	if err := base.Limit(int(limit)).Offset(int(offset)).Find(&orders).Error; err != nil {
		appLogger.L().WithError(err).Error("repo: list orders failed")
		return domain.OrderListResponse{}, err
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))
	domainOrders := make([]domain.Order, len(orders))
	for idx, ord := range orders {
		domainOrders[idx] = toDomainOrder(ord)
	}

	resp := domain.OrderListResponse{
		Data:        domainOrders,
		Total:       uint64(total),
		CurrentPage: uint64(page),
		PerPage:     uint64(limit),
		TotalInPage: uint64(len(domainOrders)),
		LastPage:    uint64(totalPages),
	}
	appLogger.L().WithField("total", resp.Total).Info("repo: list orders success")
	return resp, nil
}

// toDomainOrder converts schema to domain.
func toDomainOrder(o schema.Order) domain.Order {
	return domain.Order{
		OrderConsignmentId: o.OrderConsignmentId,
		OrderCreatedAt:     o.OrderCreatedAt,
		OrderDescription:   o.OrderDescription,
		MerchantOrderId:    o.MerchantOrderId,
		RecipientName:      o.RecipientName,
		RecipientAddress:   o.RecipientAddress,
		RecipientPhone:     o.RecipientPhone,
		OrderStatus:        domain.OrderStatus(o.OrderStatus),
		OrderAmount:        o.OrderAmount,
		TotalFee:           o.TotalFee,
		Instruction:        o.Instruction,
		OrderTypeId:        o.OrderTypeId,
		CodFee:             o.CodFee,
		PromoDiscount:      o.PromoDiscount,
		Discount:           o.Discount,
		DeliveryFee:        o.DeliveryFee,
		OrderType:          o.OrderType,
		ItemType:           o.ItemType,
		TransferStatus:     o.TransferStatus,
		Archive:            o.Archive,
		UpdatedAt:          o.UpdatedAt,
		CreatedBy:          o.CreatedBy,
		UpdatedBy:          o.UpdatedBy,
	}
}

// mapOrderDomainToSchema converts domain t0 schema.
func mapOrderDomainToSchema(params domain.Order) schema.Order {
	return schema.Order{
		OrderDescription: params.OrderDescription,
		RecipientName:    params.RecipientName,
		RecipientAddress: params.RecipientAddress,
		RecipientPhone:   params.RecipientPhone,
		OrderAmount:      params.OrderAmount,
		TotalFee:         params.TotalFee,
		MerchantOrderId:  params.MerchantOrderId,
		Instruction:      params.Instruction,
		OrderTypeId:      params.OrderTypeId,
		CodFee:           params.CodFee,
		PromoDiscount:    params.PromoDiscount,
		Discount:         params.Discount,
		DeliveryFee:      params.DeliveryFee,
		OrderStatus:      string(params.OrderStatus),
		OrderType:        params.OrderType,
		ItemType:         params.ItemType,
		TransferStatus:   params.TransferStatus,
		Archive:          params.Archive,
		CreatedBy:        params.CreatedBy,
		UpdatedBy:        params.UpdatedBy,
	}
}

// AutoMigrate ensures schema is up to date.
func (o *OrderRepository) AutoMigrate() error {
	return o.db.AutoMigrate(&schema.Order{})
}

// calculatePagination computes limit, page, and offset with defaults and bounds.
func calculatePagination(requestLimit, requestPage int64) (limit, page, offset int64) {
	const defaultPageSize int64 = 10

	limit = requestLimit
	if limit < 1 || limit > defaultPageSize {
		limit = defaultPageSize
	}

	page = requestPage
	if page < 1 {
		page = 1
	}

	offset = (page - 1) * limit
	return limit, page, offset
}

package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/domain/contract"
	"github.com/order-nest/pkg/helper"
	appLogger "github.com/order-nest/pkg/logger"
)

const (
	orderType = "Delivery"
	itemType  = "Parcel"
)

type orderUsecase struct {
	orderRepo contract.OrderRepository
}

func NewOrderUsecase(orderRepo contract.OrderRepository) contract.OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo}
}

func (o *orderUsecase) Create(ctx context.Context, params domain.CreateOrderRequest) (createResp domain.CreateOrderResponse, err error) {
	appLogger.L().WithField("created_by", params.CreatedBy).Info("create order request received")
	if params.CreatedBy == 0 {
		appLogger.L().Warn("create order failed: user id not passed")
		return createResp, errors.New("user id not passed")
	}

	// Validate input
	if err = helper.Validator.Struct(params); err != nil {
		if validationErr := convertValidationError(err, params); validationErr != nil {
			appLogger.L().WithField("created_by", params.CreatedBy).Warn("validation failed for create order")
			return createResp, validationErr
		}
		appLogger.L().WithError(err).Error("validator error for create order")
		return createResp, err
	}

	// Fee Calculation
	baseFee := 100.0
	if params.RecipientCity == 1 {
		baseFee = 60.0
	}

	var deliveryFee float64
	switch {
	case params.ItemWeight <= 0.5:
		deliveryFee = baseFee
	case params.ItemWeight <= 1.0:
		deliveryFee = baseFee + 10
	default:
		extraWeight := math.Ceil(params.ItemWeight - 1.0)
		deliveryFee = baseFee + 10 + 15*extraWeight
	}

	codFee := float64(params.AmountToCollect) * 0.01
	totalFee := deliveryFee + codFee

	order := domain.Order{
		OrderConsignmentId: uuid.New().String(),
		OrderTypeId:        1,
		PromoDiscount:      0,
		Discount:           0,
		OrderStatus:        domain.OrderStatusPending,
		OrderType:          orderType,
		ItemType:           itemType,
		TransferStatus:     1,
		Archive:            0,
		OrderCreatedAt:     time.Now(),
		OrderDescription:   params.ItemDescription,
		MerchantOrderId:    params.MerchantOrderId,
		RecipientName:      params.RecipientName,
		RecipientAddress:   params.RecipientAddress,
		RecipientPhone:     params.RecipientPhone,
		OrderAmount:        params.AmountToCollect,
		DeliveryFee:        deliveryFee,
		CodFee:             codFee,
		TotalFee:           totalFee,
		Instruction:        params.SpecialInstruction,
		CreatedBy:          params.CreatedBy,
		UpdatedBy:          params.CreatedBy,
	}

	resp, repoErr := o.orderRepo.Create(ctx, order)
	if repoErr != nil {
		appLogger.L().WithError(repoErr).Error("order repository create failed")
		return createResp, repoErr
	}
	appLogger.L().WithField("consignment_id", resp.ConsignmentId).Info("order created successfully")
	return resp, nil
}

func (o *orderUsecase) List(ctx context.Context, parameters domain.OrderListFilter) (domain.OrderListResponse, error) {
	appLogger.L().WithFields(map[string]interface{}{
		"created_by":      parameters.CreatedBy,
		"transfer_status": parameters.TransferStatus,
		"archive":         parameters.Archive,
		"limit":           parameters.Limit,
		"page":            parameters.Page,
	}).Info("list orders request received")
	resp, err := o.orderRepo.List(ctx, parameters)
	if err != nil {
		appLogger.L().WithError(err).Error("order repository list failed")
		return domain.OrderListResponse{}, err
	}
	appLogger.L().WithField("total", resp.Total).Info("list orders success")
	return resp, nil
}

// convertValidationError converts the validator's validation error into domain.ValidationError
func convertValidationError(err error, obj interface{}) *domain.ValidationError {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	errs := domain.ValidationError{ErrorMap: make(map[string][]string)}

	for _, fieldError := range validationErrors {
		jsonTag := helper.JSONTagOrFieldName(obj, fieldError.StructField())

		readable := func(j string) string {
			return strings.Join(strings.Split(j, "_"), " ")
		}

		switch tag := fieldError.ActualTag(); tag {
		case "required":
			errs.ErrorMap[jsonTag] = []string{fmt.Sprintf("The %s field is required.", readable(jsonTag))}
		case "eq":
			errs.ErrorMap[jsonTag] = []string{fmt.Sprintf("Wrong %s selected.", readable(jsonTag))}
		case "min":
			errs.ErrorMap[jsonTag] = []string{fmt.Sprintf("The %s value is too small.", readable(jsonTag))}
		case "max":
			errs.ErrorMap[jsonTag] = []string{fmt.Sprintf("The %s value is too large.", readable(jsonTag))}
		default:
			errs.ErrorMap[jsonTag] = []string{fmt.Sprintf("Field %s failed validation.", readable(jsonTag))}
		}
	}

	return &errs
}

func (o *orderUsecase) Cancel(ctx context.Context, consignmentId string, userID uint64) error {
	appLogger.L().WithFields(map[string]interface{}{
		"consignment_id": consignmentId,
		"user_id":        userID,
	}).Info("cancel order request received")
	if err := o.orderRepo.Cancel(ctx, consignmentId, userID); err != nil {
		appLogger.L().WithError(err).Error("order repository cancel failed")
		return err
	}
	appLogger.L().WithField("consignment_id", consignmentId).Info("order cancelled successfully")
	return nil
}

package domain

import (
	"time"
)

type Order struct {
	OrderConsignmentId string      `json:"order_consignment_id"`
	OrderCreatedAt     time.Time   `json:"order_created_at"`
	OrderDescription   string      `json:"order_description"`
	MerchantOrderId    string      `json:"merchant_order_id"`
	RecipientName      string      `json:"recipient_name"`
	RecipientAddress   string      `json:"recipient_address"`
	RecipientPhone     string      `json:"recipient_phone"`
	OrderAmount        int         `json:"order_amount"`
	TotalFee           float64     `json:"total_fee"`
	Instruction        string      `json:"instruction"`
	OrderTypeId        int         `json:"order_type_id"`
	CodFee             float64     `json:"cod_fee"`
	PromoDiscount      int         `json:"promo_discount"`
	Discount           int         `json:"discount"`
	DeliveryFee        float64     `json:"delivery_fee"`
	OrderStatus        OrderStatus `json:"order_status"`
	OrderType          string      `json:"order_type"`
	ItemType           string      `json:"item_type"`
	TransferStatus     uint8       `json:"transfer_status"`
	Archive            uint8       `json:"archive"`

	//audit param
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint64    `json:"created_by"`
	UpdatedBy uint64    `json:"updated_by"`
}

type CreateOrderRequest struct {
	StoreId            int     `json:"store_id,omitempty" validate:"required,eq=131172"`
	MerchantOrderId    string  `json:"merchant_order_id,omitempty"`
	RecipientName      string  `json:"recipient_name" validate:"required"`
	RecipientPhone     string  `json:"recipient_phone" validate:"required,regex_bd_phone"`
	RecipientAddress   string  `json:"recipient_address" validate:"required"`
	RecipientCity      int     `json:"recipient_city" validate:"required,eq=1"`
	RecipientZone      int     `json:"recipient_zone" validate:"required,eq=1"`
	RecipientArea      int     `json:"recipient_area" validate:"required,eq=1"`
	DeliveryType       int     `json:"delivery_type" validate:"required,eq=48"`
	ItemType           int     `json:"item_type" validate:"required,eq=2"`
	SpecialInstruction string  `json:"special_instruction"`
	ItemQuantity       int     `json:"item_quantity" validate:"required,eq=1"`
	ItemWeight         float64 `json:"item_weight" validate:"required,eq=0.5"`
	AmountToCollect    int     `json:"amount_to_collect" validate:"required,gt=0"`
	ItemDescription    string  `json:"item_description"`
	CreatedBy          uint64  `json:"created_by"`
}

type CreateOrderResponse struct {
	ConsignmentId   string  `json:"consignment_id"`
	MerchantOrderId string  `json:"merchant_order_id"`
	OrderStatus     string  `json:"order_status"`
	DeliveryFee     float64 `json:"delivery_fee"`
}

type OrderListFilter struct {
	Limit          int64  `json:"limit"`
	Page           int64  `json:"page"`
	TransferStatus uint8  `json:"transfer_status"`
	Archive        uint8  `json:"archive"`
	CreatedBy      uint64 `json:"created_by"`
}

type OrderListResponse struct {
	Data        []Order `json:"data"`
	Total       uint64  `json:"total"`
	CurrentPage uint64  `json:"current_page"`
	PerPage     uint64  `json:"per_page"`
	TotalInPage uint64  `json:"total_in_page"`
	LastPage    uint64  `json:"last_page"`
}

type OrderStatus string

const (
	OrderStatusPending  = OrderStatus("Pending")
	OrderStatusCanceled = OrderStatus("Canceled")
)

package schema

import "time"

type Order struct {
	ID                 uint64    `json:"id,omitempty" gorm:"primarykey"`
	OrderConsignmentId string    `json:"order_consignment_id" gorm:"uniqueIndex"`
	OrderCreatedAt     time.Time `json:"order_created_at"`
	OrderDescription   string    `json:"order_description"`
	MerchantOrderId    string    `json:"merchant_order_id"`
	RecipientName      string    `json:"recipient_name"`
	RecipientAddress   string    `json:"recipient_address"`
	RecipientPhone     string    `json:"recipient_phone"`
	OrderAmount        int       `json:"order_amount"`
	TotalFee           float64   `json:"total_fee"`
	Instruction        string    `json:"instruction"`
	OrderTypeId        int       `json:"order_type_id"`
	CodFee             float64   `json:"cod_fee"`
	PromoDiscount      int       `json:"promo_discount"`
	Discount           int       `json:"discount"`
	DeliveryFee        float64   `json:"delivery_fee"`
	OrderStatus        string    `json:"order_status"`
	OrderType          string    `json:"order_type"`
	ItemType           string    `json:"item_type"`
	TransferStatus     uint8     `json:"transfer_status"`
	Archive            uint8     `json:"archive"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy uint64    `json:"created_by"`
	UpdatedBy uint64    `json:"updated_by"`
}

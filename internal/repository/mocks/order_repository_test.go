package mocks

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	dialector := postgres.New(postgres.Config{Conn: db})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm with sqlmock: %v", err)
	}
	cleanup := func() { db.Close() }
	return gdb, mock, cleanup
}

func TestOrderRepository_Create(t *testing.T) {
	gdb, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := repository.NewOrderRepository(gdb)

	input := domain.Order{
		OrderDescription: "desc",
		MerchantOrderId:  "m-1",
		RecipientName:    "name",
		RecipientAddress: "addr",
		RecipientPhone:   "01700000000",
		OrderAmount:      100,
		TotalFee:         10,
		Instruction:      "note",
		OrderTypeId:      1,
		CodFee:           1,
		PromoDiscount:    0,
		Discount:         0,
		DeliveryFee:      9,
		OrderStatus:      domain.OrderStatusPending,
		OrderType:        "Delivery",
		ItemType:         "Parcel",
		TransferStatus:   1,
		Archive:          0,
		CreatedBy:        1,
		UpdatedBy:        1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), input.OrderDescription, input.MerchantOrderId, input.RecipientName, input.RecipientAddress, input.RecipientPhone, input.OrderAmount, input.TotalFee, input.Instruction, input.OrderTypeId, input.CodFee, input.PromoDiscount, input.Discount, input.DeliveryFee, string(input.OrderStatus), input.OrderType, input.ItemType, input.TransferStatus, input.Archive, sqlmock.AnyArg(), input.CreatedBy, input.UpdatedBy).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	res, err := repo.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MerchantOrderId != input.MerchantOrderId {
		t.Fatalf("unexpected merchant order id: %v", res.MerchantOrderId)
	}
}

func TestOrderRepository_List(t *testing.T) {
	gdb, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := repository.NewOrderRepository(gdb)

	params := domain.OrderListFilter{CreatedBy: 1, TransferStatus: 1, Archive: 0, Limit: 10, Page: 1}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*)`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	orderID := uuid.New().String()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"order_consignment_id", "order_created_at", "order_description", "merchant_order_id",
			"recipient_name", "recipient_address", "recipient_phone", "order_status", "order_amount",
			"total_fee", "instruction", "order_type_id", "cod_fee", "promo_discount", "discount", "delivery_fee",
			"order_type", "item_type", "transfer_status", "archive", "updated_at", "created_by", "updated_by",
		}).AddRow(
			orderID, time.Now(), "desc", "m-1", "name", "addr", "01700000000", string(domain.OrderStatusPending), 100, 10.0, "note", 1, 1.0, 0, 0, 9.0, "Delivery", "Parcel", 1, 0, time.Now(), 1, 1,
		))

	_, err := repo.List(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOrderRepository_Cancel(t *testing.T) {
	gdb, mock, cleanup := newMockDB(t)
	defer cleanup()

	repo := repository.NewOrderRepository(gdb)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE`)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	if err := repo.Cancel(context.Background(), "cons-1", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

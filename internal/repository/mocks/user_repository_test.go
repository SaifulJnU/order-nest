package mocks

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDBUser(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
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

func TestUserRepository_Create(t *testing.T) {
	gdb, mock, cleanup := newMockDBUser(t)
	defer cleanup()

	repo := repository.NewUserRepository(gdb)

	in := domain.CreateUserRequest{Username: "u1", Password: "p1"}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT`)).
		WithArgs(in.Username, in.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	_, err := repo.Create(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	gdb, mock, cleanup := newMockDBUser(t)
	defer cleanup()

	repo := repository.NewUserRepository(gdb)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs("u1", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
			AddRow(1, "u1", "p1", time.Now(), time.Now()))

	_, err := repo.GetByUsername(context.Background(), "u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	gdb, mock, cleanup := newMockDBUser(t)
	defer cleanup()

	repo := repository.NewUserRepository(gdb)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT`)).
		WithArgs(1, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
			AddRow(1, "u1", "p1", time.Now(), time.Now()))

	_, err := repo.GetByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

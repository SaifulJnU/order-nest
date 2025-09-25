package repository

import (
	"context"
	"errors"
	"time"

	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/repository/schema"
	"gorm.io/gorm"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, dUser domain.CreateUserRequest) (domain.User, error) {
	repoUser := schema.User{
		Username:  dUser.Username,
		Password:  dUser.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx := u.db.WithContext(ctx)
	if err := tx.Create(&repoUser).Error; err != nil {
		return domain.User{}, err
	}

	return mapUserSchemaToDomain(repoUser), nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, s string) (domain.User, error) {
	var repoUser schema.User
	tx := u.db.WithContext(ctx)
	if err := tx.First(&repoUser, "username = ?", s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.NotFoundError
		}
		return domain.User{}, err
	}
	return mapUserSchemaToDomain(repoUser), nil
}

func (u *UserRepository) GetByID(ctx context.Context, userId uint64) (domain.User, error) {
	var repoUser schema.User
	if err := u.db.WithContext(ctx).First(&repoUser, userId).Error; err != nil {
		return domain.User{}, err
	}
	return mapUserSchemaToDomain(repoUser), nil
}

func (u *UserRepository) AutoMigrate() error { return u.db.AutoMigrate(&schema.User{}) }

func mapUserSchemaToDomain(repoUser schema.User) domain.User {
	return domain.User{
		ID:        repoUser.ID,
		Username:  repoUser.Username,
		Password:  repoUser.Password,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt,
	}
}

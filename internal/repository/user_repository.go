package repository

import (
	"context"
	"errors"
	"time"

	"github.com/order-nest/internal/domain"
	"github.com/order-nest/internal/repository/schema"
	appLogger "github.com/order-nest/pkg/logger"
	"gorm.io/gorm"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, dUser domain.CreateUserRequest) (domain.User, error) {
	appLogger.L().WithField("username", dUser.Username).Info("repo: creating user")
	repoUser := schema.User{
		Username:  dUser.Username,
		Password:  dUser.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx := u.db.WithContext(ctx)
	if err := tx.Create(&repoUser).Error; err != nil {
		appLogger.L().WithError(err).Error("repo: user create failed")
		return domain.User{}, err
	}

	appLogger.L().WithField("user_id", repoUser.ID).Info("repo: user created")
	return mapUserSchemaToDomain(repoUser), nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, s string) (domain.User, error) {
	appLogger.L().WithField("username", s).Info("repo: get user by username")
	var repoUser schema.User
	tx := u.db.WithContext(ctx)
	if err := tx.First(&repoUser, "username = ?", s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appLogger.L().WithField("username", s).Warn("repo: user not found by username")
			return domain.User{}, domain.NotFoundError
		}
		appLogger.L().WithError(err).Error("repo: get by username failed")
		return domain.User{}, err
	}
	appLogger.L().WithField("user_id", repoUser.ID).Info("repo: user found by username")
	return mapUserSchemaToDomain(repoUser), nil
}

func (u *UserRepository) GetByID(ctx context.Context, userId uint64) (domain.User, error) {
	appLogger.L().WithField("user_id", userId).Info("repo: get user by id")
	var repoUser schema.User
	if err := u.db.WithContext(ctx).First(&repoUser, userId).Error; err != nil {
		appLogger.L().WithError(err).Error("repo: get by id failed")
		return domain.User{}, err
	}
	appLogger.L().Info("repo: user found by id")
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

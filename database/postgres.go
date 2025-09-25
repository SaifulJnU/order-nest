package database

import (
	"fmt"
	"os"

	"github.com/order-nest/config"
	appLogger "github.com/order-nest/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Postgres().Host,
		config.Postgres().User,
		config.Postgres().Password,
		config.Postgres().DbName,
		config.Postgres().Port,
		config.Postgres().SSLMode,
		config.Postgres().Timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		appLogger.L().WithError(err).Error("failed to connect to database")
		os.Exit(1)
	}
	return db
}

package database

import (
	"avto-crm-api/internal/config"
	"avto-crm-api/internal/modules/car"
	"avto-crm-api/internal/modules/client"
	"avto-crm-api/internal/modules/deal"
	"avto-crm-api/internal/modules/pipeline"
	"avto-crm-api/internal/modules/stage"
	"avto-crm-api/internal/modules/user"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormConfig)

	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к базе данных: %w", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("Ошибка при получении пула базы данных")
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	log.Println("Соединение с базой данных установлено")

	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("Ошибка с миграциями")
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Запуск миграций")

	err := db.AutoMigrate(
		&user.User{},
		&pipeline.Pipeline{},
		&car.Car{},
		&stage.Stage{},
		&client.Client{},
		&deal.Deal{},
	)

	if err != nil {
		return fmt.Errorf("Ошибка при миграции: %w", err)
	}

	log.Println("Миграции выполнены")
	return nil
}

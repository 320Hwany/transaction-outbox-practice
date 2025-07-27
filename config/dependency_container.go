package config

import (
	"gorm.io/gorm"
	"transaction-outbox-practice/application"
	"transaction-outbox-practice/controller"
	"transaction-outbox-practice/repository"
)

type Container struct {
	Config                *Config
	DB                    *gorm.DB
	OrderRepository       application.OrderRepository
	OutboxEventRepository application.OutboxEventRepository
	OrderService          *application.OrderService
	OutboxProcessor       *application.OutboxProcessor
	OrderController       *controller.OrderController
}

func NewContainer() (*Container, error) {
	cfg := LoadConfig()

	db, err := InitDB(cfg.Database.DSN())
	if err != nil {
		return nil, err
	}

	orderRepo := repository.NewOrderRepository(db)
	outboxEventRepo := repository.NewOutboxEventRepository(db)

	orderService := application.NewOrderService(db, orderRepo, outboxEventRepo)
	outboxProcessor := application.NewOutboxProcessor(outboxEventRepo, cfg.Outbox.PollingInterval, cfg.Outbox.BatchSize)
	orderController := controller.NewOrderController(orderService)

	return &Container{
		Config:                cfg,
		DB:                    db,
		OrderRepository:       orderRepo,
		OutboxEventRepository: outboxEventRepo,
		OrderService:          orderService,
		OutboxProcessor:       outboxProcessor,
		OrderController:       orderController,
	}, nil
}

func (c *Container) Close() error {
	if c.OutboxProcessor != nil {
		c.OutboxProcessor.Stop()
	}

	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}

	return nil
}

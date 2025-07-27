package application

import (
	"gorm.io/gorm"
	"transaction-outbox-practice/models"
)

type OrderRepository interface {
	Create(tx *gorm.DB, order *models.Order) error
	FindByID(id uint) (*models.Order, error)
	Update(order *models.Order) error
}

type OutboxEventRepository interface {
	Create(tx *gorm.DB, event *models.OutboxEvent) error
	FindPendingEvents(limit int) ([]models.OutboxEvent, error)
	MarkAsProcessed(event *models.OutboxEvent) error
}

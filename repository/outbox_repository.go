package repository

import (
	"gorm.io/gorm"
	"time"
	"transaction-outbox-practice/application"
	"transaction-outbox-practice/models"
)

type outboxEventRepository struct {
	db *gorm.DB
}

func NewOutboxEventRepository(db *gorm.DB) application.OutboxEventRepository {
	return &outboxEventRepository{db: db}
}

func (r *outboxEventRepository) Create(tx *gorm.DB, event *models.OutboxEvent) error {
	if tx != nil {
		return tx.Create(event).Error
	}
	return r.db.Create(event).Error
}

func (r *outboxEventRepository) FindPendingEvents(limit int) ([]models.OutboxEvent, error) {
	var events []models.OutboxEvent
	err := r.db.Where("status = ?", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *outboxEventRepository) MarkAsProcessed(event *models.OutboxEvent) error {
	now := time.Now()
	event.Status = "processed"
	event.ProcessedAt = &now
	return r.db.Save(event).Error
}

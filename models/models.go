package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	ProductName string         `gorm:"not null" json:"product_name"`
	Quantity    int            `gorm:"not null" json:"quantity"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Status      string         `gorm:"size:50;not null" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type OutboxEvent struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	AggregateID string         `gorm:"size:255;not null" json:"aggregate_id"`
	EventType   string         `gorm:"size:100;not null" json:"event_type"`
	Payload     string         `gorm:"type:text;not null" json:"payload"`
	Status      string         `gorm:"size:50;not null;default:'pending';index" json:"status"`
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`
	ProcessedAt *time.Time     `json:"processed_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

package repository

import (
	"gorm.io/gorm"
	"transaction-outbox-practice/application"
	"transaction-outbox-practice/models"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) application.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(tx *gorm.DB, order *models.Order) error {
	if tx != nil {
		return tx.Create(order).Error
	}
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

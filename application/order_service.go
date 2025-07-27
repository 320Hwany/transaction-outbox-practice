package application

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"transaction-outbox-practice/models"
)

type OrderService struct {
	db              *gorm.DB
	orderRepo       OrderRepository
	outboxEventRepo OutboxEventRepository
}

func NewOrderService(db *gorm.DB, orderRepo OrderRepository, outboxEventRepo OutboxEventRepository) *OrderService {
	return &OrderService{
		db:              db,
		orderRepo:       orderRepo,
		outboxEventRepo: outboxEventRepo,
	}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		order.Status = "pending"
		if err := s.orderRepo.Create(tx, order); err != nil {
			return err
		}

		eventPayload := map[string]interface{}{
			"order_id":     order.ID,
			"product_name": order.ProductName,
			"quantity":     order.Quantity,
			"price":        order.Price,
		}

		payloadJSON, err := json.Marshal(eventPayload)
		if err != nil {
			return err
		}

		outboxEvent := &models.OutboxEvent{
			AggregateID: fmt.Sprintf("order-%d", order.ID),
			EventType:   "OrderCreated",
			Payload:     string(payloadJSON),
			Status:      "pending",
		}

		if err := s.outboxEventRepo.Create(tx, outboxEvent); err != nil {
			return err
		}

		return nil
	})
}

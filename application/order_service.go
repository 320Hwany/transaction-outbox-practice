package application

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"transaction-outbox-practice/dto"
	"transaction-outbox-practice/models"
)

const (
	OrderStatusPending = "pending"
	EventTypeOrderCreated = "OrderCreated"
	OutboxEventStatusPending = "pending"
)

type OrderCreatedEvent struct {
	OrderID     uint    `json:"order_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}

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

func (s *OrderService) CreateOrder(req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	var response *dto.CreateOrderResponse
	
	err := s.db.Transaction(func(tx *gorm.DB) error {
		order := &models.Order{
			ProductName: req.ProductName,
			Quantity:    req.Quantity,
			Price:       req.Price,
			Status:      OrderStatusPending,
		}
		
		if err := s.orderRepo.Create(tx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		if err := s.createOrderCreatedEvent(tx, order); err != nil {
			return fmt.Errorf("failed to create outbox event: %w", err)
		}

		response = &dto.CreateOrderResponse{
			ID:          order.ID,
			ProductName: order.ProductName,
			Quantity:    order.Quantity,
			Price:       order.Price,
			Status:      order.Status,
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return response, nil
}

func (s *OrderService) createOrderCreatedEvent(tx *gorm.DB, order *models.Order) error {
	event := OrderCreatedEvent{
		OrderID:     order.ID,
		ProductName: order.ProductName,
		Quantity:    order.Quantity,
		Price:       order.Price,
	}

	payloadJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	outboxEvent := &models.OutboxEvent{
		AggregateID: fmt.Sprintf("order-%d", order.ID),
		EventType:   EventTypeOrderCreated,
		Payload:     string(payloadJSON),
		Status:      OutboxEventStatusPending,
	}

	if err := s.outboxEventRepo.Create(tx, outboxEvent); err != nil {
		return fmt.Errorf("failed to create outbox event: %w", err)
	}

	return nil
}

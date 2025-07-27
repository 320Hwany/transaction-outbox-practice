package controller

import (
	"encoding/json"
	"net/http"
	"transaction-outbox-practice/application"
	"transaction-outbox-practice/dto"
)

type OrderController struct {
	orderService *application.OrderService
}

func NewOrderController(orderService *application.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

func (c *OrderController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ProductName == "" || req.Quantity < 1 || req.Price < 0 {
		http.Error(w, "Invalid request: product_name, quantity (min 1), and price (min 0) are required", http.StatusBadRequest)
		return
	}

	response, err := c.orderService.CreateOrder(&req)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
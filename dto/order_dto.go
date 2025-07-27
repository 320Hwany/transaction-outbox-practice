package dto

type CreateOrderRequest struct {
	ProductName string  `json:"product_name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required,min=1"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

type CreateOrderResponse struct {
	ID          uint    `json:"id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
}
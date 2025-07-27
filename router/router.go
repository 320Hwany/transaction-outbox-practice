package router

import (
	"net/http"
	"transaction-outbox-practice/config"
)

func SetupRoutes(container *config.Container) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/orders", container.OrderController.CreateOrder)

	return mux
}

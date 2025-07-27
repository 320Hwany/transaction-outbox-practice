package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"transaction-outbox-practice/config"
	"transaction-outbox-practice/models"
)

func main() {
	container, err := config.NewContainer()
	if err != nil {
		log.Fatal("Failed to initialize DI container:", err)
	}
	defer container.Close()

	go container.OutboxProcessor.Start()

	order := &models.Order{
		ProductName: "test product name",
		Quantity:    3,
		Price:       1000,
	}

	fmt.Println("Creating order...")
	if err := container.OrderService.CreateOrder(order); err != nil {
		log.Fatal("Failed to create order:", err)
	}

	fmt.Printf("Order created successfully! ID: %d\n", order.ID)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\nShutting down...")
}

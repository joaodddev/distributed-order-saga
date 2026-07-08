package main

import (
	"log"
	"os"

	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/usecase"
	httpinfra "github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/http"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/http/handler"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/persistence/mysql"
)

func main() {
	db, err := mysql.NewConnection(mysql.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Database: getEnv("DB_NAME", "order_db"),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	orderRepository := mysql.NewOrderRepository(db)
	createOrderUseCase := usecase.NewCreateOrder(orderRepository)
	orderHandler := handler.NewOrderHandler(createOrderUseCase)

	router := httpinfra.NewRouter(orderHandler)

	port := getEnv("PORT", "8080")
	log.Printf("order-service listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

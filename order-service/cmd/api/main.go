package main

import (
	"context"
	"log"
	"os"

	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/usecase"
	httpinfra "github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/http"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/http/handler"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/messaging"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/observability"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/persistence/mysql"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/outbox"
)

func main() {
	ctx := context.Background()

	shutdown, err := observability.InitTracer(ctx, "order-service", getEnv("OTEL_COLLECTOR_ENDPOINT", "localhost:4317"))
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer shutdown(ctx)

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

	kafkaBrokers := []string{getEnv("KAFKA_BROKERS", "localhost:9092")}

	producer := messaging.NewProducer(kafkaBrokers)
	defer producer.Close()

	outboxRepository := mysql.NewOutboxRepository(db)
	relay := outbox.NewRelay(outboxRepository, producer)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go relay.Start(ctx)

	orderRepository := mysql.NewOrderRepository(db)

	createOrderUseCase := usecase.NewCreateOrder(orderRepository)
	orderHandler := handler.NewOrderHandler(createOrderUseCase)

	cancelOrderUseCase := usecase.NewCancelOrder(orderRepository)
	compensationConsumer := messaging.NewConsumer(kafkaBrokers, "order-service-compensation-group", cancelOrderUseCase)
	defer compensationConsumer.Close()
	go compensationConsumer.Start(ctx)

	confirmOrderUseCase := usecase.NewConfirmOrder(orderRepository)
	confirmationConsumer := messaging.NewConfirmationConsumer(kafkaBrokers, "order-service-confirmation-group", confirmOrderUseCase)
	defer confirmationConsumer.Close()
	go confirmationConsumer.Start(ctx)

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

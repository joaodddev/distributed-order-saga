package main

import (
	"context"
	"log"
	"os"

	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/application/usecase"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/infrastructure/messaging"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/infrastructure/persistence/mysql"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/outbox"
)

func main() {
	db, err := mysql.NewConnection(mysql.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3307"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Database: getEnv("DB_NAME", "inventory_db"),
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go relay.Start(ctx)

	stockRepository := mysql.NewStockRepository(db)
	reserveStockUseCase := usecase.NewReserveStock(stockRepository)

	consumer := messaging.NewConsumer(kafkaBrokers, "inventory-service-group", reserveStockUseCase)
	defer consumer.Close()

	log.Println("inventory-service worker started, consuming payment.reserved")
	consumer.Start(ctx)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

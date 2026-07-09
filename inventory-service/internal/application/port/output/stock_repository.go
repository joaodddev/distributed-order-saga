package output

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/domain"
)

type OutboxEvent struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
}

type StockRepository interface {
	SaveWithOutboxEvent(ctx context.Context, reservation *domain.StockReservation, event OutboxEvent) error
}

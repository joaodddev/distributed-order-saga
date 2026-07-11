package output

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/domain"
)

type OutboxEvent struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
}

type OrderRepository interface {
	SaveWithOutboxEvent(ctx context.Context, order *domain.Order, event OutboxEvent) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	CancelWithOutboxEvent(ctx context.Context, orderID uuid.UUID, event OutboxEvent) error
}

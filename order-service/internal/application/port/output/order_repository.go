package output

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/domain"
)

// OutboxEvent representa a linha gravada na tabela outbox_events,
// sempre na MESMA transação SQL que o insert do pedido.
type OutboxEvent struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
}

type OrderRepository interface {
	// SaveWithOutboxEvent persiste o pedido e o evento de outbox atomicamente.
	// É essa atomicidade que garante o "transactional" do Transactional Outbox.
	SaveWithOutboxEvent(ctx context.Context, order *domain.Order, event OutboxEvent) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
}

CancelWithOutboxEvent(ctx context.Context, orderID uuid.UUID, event OutboxEvent) error

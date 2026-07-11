package output

import (
	"context"
	"fmt"
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

func (r *OrderRepository) ConfirmWithOutboxEvent(ctx context.Context, orderID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE orders SET status = ? WHERE id = ?`,
		domain.OrderStatusConfirmed, orderID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to confirm order: %w", err)
	}
	return nil
}

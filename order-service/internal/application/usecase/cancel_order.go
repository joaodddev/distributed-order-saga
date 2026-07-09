package usecase

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/domain"
)

// PaymentRefundedPayload espelha o payload do evento payment.refunded.
type PaymentRefundedPayload struct {
	OrderID string `json:"orderId"`
	Reason  string `json:"reason"`
}

type CancelOrder struct {
	repository output.OrderRepository
}

func NewCancelOrder(repository output.OrderRepository) *CancelOrder {
	return &CancelOrder{repository: repository}
}

func (uc *CancelOrder) Execute(ctx context.Context, payload PaymentRefundedPayload, sagaID, correlationID string) error {
	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return err
	}

	event := domain.NewOrderCancelledEvent(orderID.String(), sagaID, correlationID, payload.Reason)
	eventPayload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outboxEvent := output.OutboxEvent{
		ID:          uuid.New(),
		AggregateID: orderID,
		EventType:   event.EventType,
		Payload:     eventPayload,
		CreatedAt:   event.OccurredAt,
	}

	return uc.repository.CancelWithOutboxEvent(ctx, orderID, outboxEvent)
}

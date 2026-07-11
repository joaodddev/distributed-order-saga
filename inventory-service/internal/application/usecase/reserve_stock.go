package usecase

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/domain"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/infrastructure/observability"
)

type PaymentReservedPayload struct {
	OrderID    string `json:"orderId"`
	CustomerID string `json:"customerId"`
}

type ReserveStock struct {
	repository output.StockRepository
}

func NewReserveStock(repository output.StockRepository) *ReserveStock {
	return &ReserveStock{repository: repository}
}

func (uc *ReserveStock) Execute(ctx context.Context, payload PaymentReservedPayload, sagaID, correlationID string) error {
	tracer := observability.Tracer("inventory-service")
	ctx, span := tracer.Start(ctx, "ReserveStock.Execute")
	defer span.End()

	span.SetAttributes(
		attribute.String("saga.id", sagaID),
		attribute.String("order.id", payload.OrderID),
	)

	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		span.RecordError(err)
		return err
	}
	customerID, err := uuid.Parse(payload.CustomerID)
	if err != nil {
		span.RecordError(err)
		return err
	}

	available := orderID.String()[len(orderID.String())-1] != '0'

	reservation := domain.Reserve(orderID, customerID, available)
	span.SetAttributes(attribute.String("reservation.status", string(reservation.Status)))

	reason := ""
	if !reservation.Reserved() {
		reason = "insufficient stock"
	}

	event := domain.NewInventoryEvent(reservation, sagaID, correlationID, reason)
	eventPayload, err := json.Marshal(event)
	if err != nil {
		span.RecordError(err)
		return err
	}

	outboxEvent := output.OutboxEvent{
		ID:          uuid.New(),
		AggregateID: reservation.ID,
		EventType:   event.EventType,
		Payload:     eventPayload,
		CreatedAt:   reservation.CreatedAt,
	}

	if err := uc.repository.SaveWithOutboxEvent(ctx, reservation, outboxEvent); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}

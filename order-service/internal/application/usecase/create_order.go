package usecase

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/input"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/domain"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/infrastructure/observability"
)

type CreateOrder struct {
	repository output.OrderRepository
}

func NewCreateOrder(repository output.OrderRepository) *CreateOrder {
	return &CreateOrder{repository: repository}
}

func (uc *CreateOrder) Execute(ctx context.Context, in input.CreateOrderInput) (*input.CreateOrderOutput, error) {
	tracer := observability.Tracer("order-service")
	ctx, span := tracer.Start(ctx, "CreateOrder.Execute")
	defer span.End()

	items := make([]domain.OrderItem, len(in.Items))
	for i, item := range in.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	order, err := domain.NewOrder(in.CustomerID, items)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(
		attribute.String("saga.id", order.ID.String()),
		attribute.String("order.id", order.ID.String()),
		attribute.Float64("order.total_amount", order.TotalAmount),
	)

	event := domain.NewOrderCreatedEvent(order)
	payload, err := json.Marshal(event)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	outboxEvent := output.OutboxEvent{
		ID:          uuid.New(),
		AggregateID: order.ID,
		EventType:   event.EventType,
		Payload:     payload,
		CreatedAt:   order.CreatedAt,
	}

	if err := uc.repository.SaveWithOutboxEvent(ctx, order, outboxEvent); err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &input.CreateOrderOutput{
		OrderID: order.ID,
		Status:  string(order.Status),
	}, nil
}

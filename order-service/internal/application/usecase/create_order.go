package usecase

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/input"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/domain"
)

type CreateOrder struct {
	repository output.OrderRepository
}

func NewCreateOrder(repository output.OrderRepository) *CreateOrder {
	return &CreateOrder{repository: repository}
}

func (uc *CreateOrder) Execute(ctx context.Context, in input.CreateOrderInput) (*input.CreateOrderOutput, error) {
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
		return nil, err
	}

	event := domain.NewOrderCreatedEvent(order)
	payload, err := json.Marshal(event)
	if err != nil {
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
		return nil, err
	}

	return &input.CreateOrderOutput{
		OrderID: order.ID,
		Status:  string(order.Status),
	}, nil
}

package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/output"
)

type OrderConfirmedPayload struct {
	OrderID    string `json:"orderId"`
	CustomerID string `json:"customerId"`
}

type ConfirmOrder struct {
	repository output.OrderRepository
}

func NewConfirmOrder(repository output.OrderRepository) *ConfirmOrder {
	return &ConfirmOrder{repository: repository}
}

func (uc *ConfirmOrder) Execute(ctx context.Context, payload OrderConfirmedPayload) error {
	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return err
	}
	return uc.repository.ConfirmWithOutboxEvent(ctx, orderID)
}

package input

import (
	"context"

	"github.com/google/uuid"
)

type CreateOrderItemInput struct {
	ProductID uuid.UUID
	Quantity  int
	UnitPrice float64
}

type CreateOrderInput struct {
	CustomerID uuid.UUID
	Items      []CreateOrderItemInput
}

type CreateOrderOutput struct {
	OrderID uuid.UUID
	Status  string
}

// CreateOrderUseCase é a porta que o handler HTTP vai depender —
// nunca da implementação concreta, só dessa interface.
type CreateOrderUseCase interface {
	Execute(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error)
}

package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

var (
	ErrEmptyItems      = errors.New("order must contain at least one item")
	ErrInvalidQuantity = errors.New("item quantity must be greater than zero")
)

type OrderItem struct {
	ProductID uuid.UUID
	Quantity  int
	UnitPrice float64
}

type Order struct {
	ID          uuid.UUID
	CustomerID  uuid.UUID
	Items       []OrderItem
	TotalAmount float64
	Status      OrderStatus
	CreatedAt   time.Time
}

func NewOrder(customerID uuid.UUID, items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyItems
	}

	var total float64
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
		total += item.UnitPrice * float64(item.Quantity)
	}

	return &Order{
		ID:          uuid.New(),
		CustomerID:  customerID,
		Items:       items,
		TotalAmount: total,
		Status:      OrderStatusPending,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (o *Order) Confirm() {
	o.Status = OrderStatusConfirmed
}

func (o *Order) Cancel() {
	o.Status = OrderStatusCancelled
}

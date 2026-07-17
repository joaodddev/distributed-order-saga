package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/input"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/usecase"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/domain"
)

// fakeOrderRepository é um double em memória — não usa mock framework,
// já que a interface OrderRepository é pequena o suficiente pra isso ser
// mais simples e legível do que introduzir gomock/testify/mock.
type fakeOrderRepository struct {
	savedOrder *domain.Order
	savedEvent output.OutboxEvent
	saveErr    error
}

func (f *fakeOrderRepository) SaveWithOutboxEvent(ctx context.Context, order *domain.Order, event output.OutboxEvent) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.savedOrder = order
	f.savedEvent = event
	return nil
}

func (f *fakeOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	return f.savedOrder, nil
}

func (f *fakeOrderRepository) CancelWithOutboxEvent(ctx context.Context, orderID uuid.UUID, event output.OutboxEvent) error {
	return nil
}

func (f *fakeOrderRepository) ConfirmWithOutboxEvent(ctx context.Context, orderID uuid.UUID) error {
	return nil
}

func TestCreateOrder_Execute_Success(t *testing.T) {
	repo := &fakeOrderRepository{}
	uc := usecase.NewCreateOrder(repo)

	out, err := uc.Execute(context.Background(), input.CreateOrderInput{
		CustomerID: uuid.New(),
		Items: []input.CreateOrderItemInput{
			{ProductID: uuid.New(), Quantity: 2, UnitPrice: 49.90},
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out.Status != string(domain.OrderStatusPending) {
		t.Errorf("expected status PENDING, got %s", out.Status)
	}
	if repo.savedOrder == nil {
		t.Fatal("expected order to be saved, but it wasn't")
	}
	if repo.savedEvent.EventType != "order.created" {
		t.Errorf("expected outbox event type order.created, got %s", repo.savedEvent.EventType)
	}
}

func TestCreateOrder_Execute_EmptyItems_ReturnsError(t *testing.T) {
	repo := &fakeOrderRepository{}
	uc := usecase.NewCreateOrder(repo)

	_, err := uc.Execute(context.Background(), input.CreateOrderInput{
		CustomerID: uuid.New(),
		Items:      []input.CreateOrderItemInput{},
	})

	if err == nil {
		t.Fatal("expected error for empty items, got nil")
	}
}

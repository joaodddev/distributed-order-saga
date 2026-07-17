package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/application/usecase"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/domain"
)

type fakeStockRepository struct {
	savedReservation *domain.StockReservation
	savedEvent       output.OutboxEvent
}

func (f *fakeStockRepository) SaveWithOutboxEvent(ctx context.Context, reservation *domain.StockReservation, event output.OutboxEvent) error {
	f.savedReservation = reservation
	f.savedEvent = event
	return nil
}

func TestReserveStock_Execute_SavesReservationAndEvent(t *testing.T) {
	repo := &fakeStockRepository{}
	uc := usecase.NewReserveStock(repo)

	payload := usecase.PaymentReservedPayload{
		OrderID:    uuid.New().String(),
		CustomerID: uuid.New().String(),
	}

	err := uc.Execute(context.Background(), payload, "saga-1", "corr-1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.savedReservation == nil {
		t.Fatal("expected reservation to be saved, but it wasn't")
	}
	if repo.savedEvent.EventType != "inventory.reserved" && repo.savedEvent.EventType != "inventory.failed" {
		t.Errorf("unexpected event type: %s", repo.savedEvent.EventType)
	}
}

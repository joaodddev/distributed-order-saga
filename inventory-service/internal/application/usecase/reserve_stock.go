package usecase

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/domain"
)

// PaymentReservedPayload espelha o payload do evento payment.reserved
// publicado pelo payment-service (mesmo contrato definido em contracts/events).
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
	orderID, err := uuid.Parse(payload.OrderID)
	if err != nil {
		return err
	}
	customerID, err := uuid.Parse(payload.CustomerID)
	if err != nil {
		return err
	}

	// Simulação simplificada: em produção isso consultaria uma tabela real
	// de stock_items. Aqui, pedidos de um customer terminado em "0" falham
	// de propósito, só pra ter um cenário de compensação demonstrável no vídeo.
	available := orderID.String()[len(orderID.String())-1] != '0'

	reservation := domain.Reserve(orderID, customerID, available)

	reason := ""
	if !reservation.Reserved() {
		reason = "insufficient stock"
	}

	event := domain.NewInventoryEvent(reservation, sagaID, correlationID, reason)
	eventPayload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	outboxEvent := output.OutboxEvent{
		ID:          uuid.New(),
		AggregateID: reservation.ID,
		EventType:   event.EventType,
		Payload:     eventPayload,
		CreatedAt:   reservation.CreatedAt,
	}

	return uc.repository.SaveWithOutboxEvent(ctx, reservation, outboxEvent)
}

package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReservationStatus string

const (
	ReservationStatusReserved ReservationStatus = "RESERVED"
	ReservationStatusFailed   ReservationStatus = "FAILED"
	ReservationStatusReleased ReservationStatus = "RELEASED"
)

type StockReservation struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	Status     ReservationStatus
	CreatedAt  time.Time
}

// Reserve simula a checagem de estoque. Numa implementação real isso
// consultaria uma tabela stock_items por produto; aqui simplificamos
// pra manter o foco no padrão da saga em vez de regra de negócio de estoque.
func Reserve(orderID, customerID uuid.UUID, available bool) *StockReservation {
	status := ReservationStatusReserved
	if !available {
		status = ReservationStatusFailed
	}

	return &StockReservation{
		ID:         uuid.New(),
		OrderID:    orderID,
		CustomerID: customerID,
		Status:     status,
		CreatedAt:  time.Now().UTC(),
	}
}

func (s *StockReservation) Reserved() bool {
	return s.Status == ReservationStatusReserved
}

package domain

import "time"

type InventoryReservedEvent struct {
	EventType     string                `json:"eventType"`
	Version       int                   `json:"version"`
	SagaID        string                `json:"sagaId"`
	CorrelationID string                `json:"correlationId"`
	OccurredAt    time.Time             `json:"occurredAt"`
	Payload       InventoryEventPayload `json:"payload"`
}

type InventoryEventPayload struct {
	ReservationID string `json:"reservationId"`
	OrderID       string `json:"orderId"`
	CustomerID    string `json:"customerId"`
	Reason        string `json:"reason,omitempty"`
}

func NewInventoryEvent(reservation *StockReservation, sagaID, correlationID, reason string) InventoryReservedEvent {
	eventType := "inventory.reserved"
	if !reservation.Reserved() {
		eventType = "inventory.failed"
	}

	return InventoryReservedEvent{
		EventType:     eventType,
		Version:       1,
		SagaID:        sagaID,
		CorrelationID: correlationID,
		OccurredAt:    reservation.CreatedAt,
		Payload: InventoryEventPayload{
			ReservationID: reservation.ID.String(),
			OrderID:       reservation.OrderID.String(),
			CustomerID:    reservation.CustomerID.String(),
			Reason:        reason,
		},
	}
}

package domain

import "time"

type OrderCreatedEvent struct {
	EventType     string              `json:"eventType"`
	Version       int                 `json:"version"`
	SagaID        string              `json:"sagaId"`
	CorrelationID string              `json:"correlationId"`
	OccurredAt    time.Time           `json:"occurredAt"`
	Payload       OrderCreatedPayload `json:"payload"`
}

type OrderCreatedPayload struct {
	OrderID     string             `json:"orderId"`
	CustomerID  string             `json:"customerId"`
	Items       []OrderCreatedItem `json:"items"`
	TotalAmount float64            `json:"totalAmount"`
}

type OrderCreatedItem struct {
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
}

func NewOrderCreatedEvent(order *Order) OrderCreatedEvent {
	items := make([]OrderCreatedItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = OrderCreatedItem{
			ProductID: item.ProductID.String(),
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		}
	}

	// Saga e correlation nascem juntos com o pedido: todo evento subsequente
	// na cadeia vai carregar esse mesmo SagaID pra permitir rastreamento no Jaeger.
	return OrderCreatedEvent{
		EventType:     "order.created",
		Version:       1,
		SagaID:        order.ID.String(),
		CorrelationID: order.ID.String(),
		OccurredAt:    order.CreatedAt,
		Payload: OrderCreatedPayload{
			OrderID:     order.ID.String(),
			CustomerID:  order.CustomerID.String(),
			Items:       items,
			TotalAmount: order.TotalAmount,
		},
	}
}

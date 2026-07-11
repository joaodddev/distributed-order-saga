package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/usecase"
)

type orderConfirmedEnvelope struct {
	Payload usecase.OrderConfirmedPayload `json:"payload"`
}

type ConfirmationConsumer struct {
	reader  *kafka.Reader
	useCase *usecase.ConfirmOrder
}

func NewConfirmationConsumer(brokers []string, groupID string, useCase *usecase.ConfirmOrder) *ConfirmationConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   "order.confirmed",
	})
	return &ConfirmationConsumer{reader: reader, useCase: useCase}
}

func (c *ConfirmationConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("confirmation consumer: failed to read message: %v", err)
			continue
		}

		var envelope orderConfirmedEnvelope
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			log.Printf("confirmation consumer: failed to unmarshal message: %v", err)
			continue
		}

		if err := c.useCase.Execute(ctx, envelope.Payload); err != nil {
			log.Printf("confirmation consumer: failed to confirm order: %v", err)
			continue
		}

		log.Printf("[order.confirmed] order %s confirmed", envelope.Payload.OrderID)
	}
}

func (c *ConfirmationConsumer) Close() error {
	return c.reader.Close()
}

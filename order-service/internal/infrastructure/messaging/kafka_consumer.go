package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/usecase"
)

type paymentRefundedEnvelope struct {
	SagaID        string                         `json:"sagaId"`
	CorrelationID string                         `json:"correlationId"`
	Payload       usecase.PaymentRefundedPayload `json:"payload"`
}

type Consumer struct {
	reader  *kafka.Reader
	useCase *usecase.CancelOrder
}

func NewConsumer(brokers []string, groupID string, useCase *usecase.CancelOrder) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   "payment.refunded",
	})

	return &Consumer{reader: reader, useCase: useCase}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("order consumer: failed to read message: %v", err)
			continue
		}

		var envelope paymentRefundedEnvelope
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			log.Printf("order consumer: failed to unmarshal message: %v", err)
			continue
		}

		if err := c.useCase.Execute(ctx, envelope.Payload, envelope.SagaID, envelope.CorrelationID); err != nil {
			log.Printf("order consumer: failed to process payment.refunded: %v", err)
			continue
		}

		log.Printf("[payment.refunded] cancelled order %s", envelope.Payload.OrderID)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

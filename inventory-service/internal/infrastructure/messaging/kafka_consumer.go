package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/application/usecase"
)

type paymentReservedEnvelope struct {
	SagaID        string                         `json:"sagaId"`
	CorrelationID string                         `json:"correlationId"`
	Payload       usecase.PaymentReservedPayload `json:"payload"`
}

type Consumer struct {
	reader  *kafka.Reader
	useCase *usecase.ReserveStock
}

func NewConsumer(brokers []string, groupID string, useCase *usecase.ReserveStock) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   "payment.reserved",
	})

	return &Consumer{reader: reader, useCase: useCase}
}

// Start roda em loop bloqueante. Se ReserveStock falhar, a mensagem NÃO é
// commitada (ReadMessage/FetchMessage + CommitMessages seria o padrão pra
// controle manual; aqui usamos ReadMessage que já avança o offset — então
// falha aqui significa perda do evento. Isso é aceitável nessa fase do
// projeto e vira ponto de melhoria quando adicionarmos DLQ.
func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("inventory consumer: failed to read message: %v", err)
			continue
		}

		var envelope paymentReservedEnvelope
		if err := json.Unmarshal(msg.Value, &envelope); err != nil {
			log.Printf("inventory consumer: failed to unmarshal message: %v", err)
			continue
		}

		if err := c.useCase.Execute(ctx, envelope.Payload, envelope.SagaID, envelope.CorrelationID); err != nil {
			log.Printf("inventory consumer: failed to process payment.reserved: %v", err)
			continue
		}

		log.Printf("[payment.reserved] processed order %s", envelope.Payload.OrderID)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

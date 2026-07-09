package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/inventory-service/internal/domain"
)

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) SaveWithOutboxEvent(ctx context.Context, reservation *domain.StockReservation, event output.OutboxEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO stock_reservations (id, order_id, customer_id, status, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		reservation.ID.String(), reservation.OrderID.String(), reservation.CustomerID.String(),
		reservation.Status, reservation.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert stock reservation: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO outbox_events (id, aggregate_id, event_type, payload, published, created_at)
		 VALUES (?, ?, ?, ?, FALSE, ?)`,
		event.ID.String(), event.AggregateID.String(), event.EventType, event.Payload, event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert outbox event: %w", err)
	}

	return tx.Commit()
}

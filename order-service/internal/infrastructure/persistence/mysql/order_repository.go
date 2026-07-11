package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/application/port/output"
	"github.com/joaodddev/distributed-order-saga/order-service/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) SaveWithOutboxEvent(ctx context.Context, order *domain.Order, event output.OutboxEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO orders (id, customer_id, total_amount, status, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		order.ID.String(), order.CustomerID.String(), order.TotalAmount, order.Status, order.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO order_items (id, order_id, product_id, quantity, unit_price)
			 VALUES (?, ?, ?, ?, ?)`,
			uuid.New().String(), order.ID.String(), item.ProductID.String(), item.Quantity, item.UnitPrice,
		)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
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

func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, customer_id, total_amount, status, created_at
		 FROM orders WHERE id = ?`, id.String(),
	)

	var order domain.Order
	var orderID, customerID string
	if err := row.Scan(&orderID, &customerID, &order.TotalAmount, &order.Status, &order.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found: %s", id)
		}
		return nil, fmt.Errorf("failed to query order: %w", err)
	}

	order.ID = uuid.MustParse(orderID)
	order.CustomerID = uuid.MustParse(customerID)

	return &order, nil
}

func (r *OrderRepository) CancelWithOutboxEvent(ctx context.Context, orderID uuid.UUID, event output.OutboxEvent) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE orders SET status = ? WHERE id = ?`,
		domain.OrderStatusCancelled, orderID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
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

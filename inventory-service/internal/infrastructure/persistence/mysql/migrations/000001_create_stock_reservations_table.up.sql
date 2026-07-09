CREATE TABLE stock_reservations (
    id CHAR(36) PRIMARY KEY,
    order_id CHAR(36) NOT NULL,
    customer_id CHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,

    INDEX idx_stock_reservations_order_id (order_id)
) ENGINE=InnoDB;

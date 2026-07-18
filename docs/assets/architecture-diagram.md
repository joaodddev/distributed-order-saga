# Architecture diagram

This file can be used as a lightweight visual reference for the project architecture.

```text
+-------------------+       +-------------------+       +-------------------+
| order-service     | ----> | payment-service   | ----> | inventory-service |
| (Go + Gin)       |       | (Ruby + Roda)    |       | (Go worker)      |
+-------------------+       +-------------------+       +-------------------+
         |                             |                           |
         |                             |                           |
         v                             v                           v
+-------------------+       +-------------------+       +-------------------+
| Kafka topics      |       | PostgreSQL        |       | MySQL            |
| order.created     |       | payments/outbox  |       | stock/outbox    |
+-------------------+       +-------------------+       +-------------------+
                                 |
                                 v
                        +-------------------+
                        | notification-service |
                        | (Ruby + Racecar)   |
                        +-------------------+
```

Suggested screenshot sources:
- Jaeger UI trace view for a saga execution.
- Docker Compose service overview showing Kafka, MySQL, PostgreSQL, Redis and Jaeger.

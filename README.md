# distributed-order-saga

Distributed order processing system implementing the Choreography-based Saga pattern with Transactional Outbox, built with Go and Ruby and powered by Kafka.

This monorepo contains four independent services that collaborate through events to process orders end-to-end:

- order-service: creates orders, starts the saga, and manages confirmation/cancellation.
- payment-service: reserves and refunds payments.
- inventory-service: reserves and releases stock.
- notification-service: publishes order confirmation notifications idempotently.

## Architecture

The flow is choreography-based and event-driven:

```text
order.created -> payment.reserved -> inventory.reserved -> order.confirmed
                                    -> inventory.failed -> payment.refunded -> order.cancelled
```

Every event carries a sagaId and a correlationId created in the initial order.created event and propagated throughout the workflow.

### Visual references

- Architecture reference: [docs/assets/architecture-diagram.md](docs/assets/architecture-diagram.md)
- Demo video guide: [docs/assets/demo-video.md](docs/assets/demo-video.md)

A good screenshot to include in the repository README is a Jaeger trace for a completed saga or a Docker Compose overview showing the event-driven infrastructure.

## Services and stack

| Service | Stack | Role |
| --- | --- | --- |
| order-service | Go + Gin + MySQL | Creates orders and coordinates the saga lifecycle |
| payment-service | Ruby + Roda + Racecar + PostgreSQL | Reserves and refunds payments |
| inventory-service | Go + MySQL | Reserves stock and processes compensation |
| notification-service | Ruby + Racecar + Redis | Emits confirmation notifications idempotently |

## Technical decisions

- Transactional Outbox is used in the Go services and in the payment-service to reliably publish domain events without losing messages.
- Choreography was chosen over orchestration because the services already share explicit domain events and the workflow stays decoupled and easy to evolve.
- The notification-service uses Redis-based idempotency instead of outbox because it only needs to ensure a notification is emitted once per saga and does not participate in the transactional write path of the core business process.

## Local setup

### 1. Start infrastructure

```bash
cd ~/colabs/distributed-order-saga
docker-compose up -d mysql mysql-inventory postgres redis kafka jaeger loki grafana
```

If a container was removed and you need to recreate the schemas, run:

```bash
cd order-service && make migrate-up
cd ../inventory-service && make migrate-up
```

### 2. Start each service

```bash
cd order-service && go run ./cmd/api
cd ../inventory-service && go run ./cmd/worker
cd ../payment-service && bundle exec rackup -p 4567
cd ../notification-service && bundle exec racecar
```

### 3. Run tests

```bash
cd order-service && go test ./...
cd ../inventory-service && go test ./...
cd ../payment-service && bundle exec rspec
cd ../notification-service && bundle exec rspec
```

## Observability

The services export tracing to Jaeger through OTLP:

- Go services: localhost:4317
- Ruby services: localhost:4318

Open the Jaeger UI at http://localhost:16686 and search by the same sagaId to inspect the full distributed trace.

## CI

Continuous integration is configured in [.github/workflows/ci.yml](.github/workflows/ci.yml) and runs build and test checks for each service independently.

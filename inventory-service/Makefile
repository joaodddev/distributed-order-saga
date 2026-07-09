.PHONY: up down reset-db migrate-up migrate-down run

up:
	cd .. && docker-compose up -d mysql-inventory

down:
	cd .. && docker-compose down

reset-db:
	cd .. && docker-compose down -v
	cd .. && docker-compose up -d mysql-inventory
	sleep 3
	$(MAKE) migrate-up

migrate-up:
	migrate -database "mysql://root:root@tcp(localhost:3307)/inventory_db" \
		-path internal/infrastructure/persistence/mysql/migrations up

migrate-down:
	migrate -database "mysql://root:root@tcp(localhost:3307)/inventory_db" \
		-path internal/infrastructure/persistence/mysql/migrations down

run:
	go run cmd/worker/main.go

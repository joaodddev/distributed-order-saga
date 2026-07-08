.PHONY: up down reset-db migrate-up migrate-down run

up:
	cd .. && docker-compose up -d mysql

down:
	cd .. && docker-compose down

reset-db:
	cd .. && docker-compose down -v
	cd .. && docker-compose up -d mysql
	sleep 3
	$(MAKE) migrate-up

migrate-up:
	migrate -database "mysql://root:root@tcp(localhost:3306)/order_db" \
		-path internal/infrastructure/persistence/migrations up

migrate-down:
	migrate -database "mysql://root:root@tcp(localhost:3306)/order_db" \
		-path internal/infrastructure/persistence/migrations down

run:
	go run cmd/api/main.go

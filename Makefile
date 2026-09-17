start:
	docker start order-postgres
	go run ./cmd/internal/models

stop:
	docker stop order-postgres

test:
	go test ./...

build:
	go build ./...

run:
	go run ./cmd/internal/models

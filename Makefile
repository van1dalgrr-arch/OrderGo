start:
	docker start order-postgres
	go run ./cmd/internal/models

stop:
	docker stop order-postgres

run:
	go run ./cmd/internal/models

test:
	go test ./...

build:
	go build ./...

up:
	docker compose --env-file .env -f deploy/docker-compose.yml up --build -d

down:
	docker compose --env-file .env -f deploy/docker-compose.yml down

logs:
	docker compose --env-file .env -f deploy/docker-compose.yml logs -f


stat:
	git status

base:
	git add . && git commit -m "$(msg)" && git push

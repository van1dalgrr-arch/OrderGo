.PHONY: run test build up down logs stat base form db

run:
	@go run ./cmd/api

test:
	@go test ./...

build:
	@go build ./...

up:
	@docker compose --env-file .env -f deploy/docker-compose.yml up --build -d

down:
	@docker compose --env-file .env -f deploy/docker-compose.yml down

logs:
	@docker compose --env-file .env -f deploy/docker-compose.yml logs -f

stat:
	@git status

base:
	@git add . && git commit -m "$(msg)" && git push

form:
	@gofmt -w .

db:
	@docker compose --env-file .env -f deploy/docker-compose.yml exec postgres psql -U postgres -d orderapi

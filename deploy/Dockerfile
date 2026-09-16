FROM golang:1.27.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order-api ./cmd/internal/models

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/order-api .

EXPOSE 8080

CMD ["./order-api"]
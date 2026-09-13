🧾 Order API

Backend API для управления заказами, написанный на Go + Gin.

Проект развивается от простого in-memory API к полноценной системе заказов и оплаты с PostgreSQL и Telegram-ботом.

⸻

🛠 Stack

Technology	Purpose
🐹 Go	Backend
🌐 Gin	HTTP framework
🗺️ Map	Temporary in-memory storage
🐳 Docker	Containerization
💳 Stripe	Planned payment integration
🐘 PostgreSQL	Planned database
🤖 Telegram Bot	Planned interface

⸻

📊 Current Progress

███████████████░░░░░░░░ 50%

API

* HTTP server
* Health check
* Create order
* Get order by ID
* Delete order
* HTTP error handling
* Payment flow
* Order status management
* PostgreSQL
* Telegram bot

⸻

🏗️ Current Architecture

flowchart TD
Client[Client]
Gin[Gin HTTP Server]
Create[Create Order]
Get[Get Order]
Delete[Delete Order]
Storage[(In-Memory Map)]
Client --> Gin
Gin --> Create
Gin --> Get
Gin --> Delete
Create --> Storage
Get --> Storage
Delete --> Storage

⸻

📦 Order Lifecycle

Планируемая логика заказа:

stateDiagram-v2
[*] --> New
New --> Paid
New --> Cancelled
Paid --> [*]
Cancelled --> [*]

⸻

🔌 API

Health Check

GET /health

Response:

{
"status": "ok"
}

Create Order

POST /orders

Request:

{
"user_id": "123",
"amount": 1999,
"status": "new"
}

Response:

{
"id": "1",
"user_id": "123",
"amount": 1999,
"status": "new"
}

amount хранится как int в минимальных денежных единицах. float для денег не используется из-за проблем с точностью представления.

Get Order

GET /orders/:id

Example:

GET /orders/1

Delete Order

DELETE /orders/:id

Response:

{
"message": "order has been deleted"
}

⸻

🚀 Run Locally

Clone the repository:

git clone <repository-url>
cd order-api

Install dependencies:

go mod download

Run:

go run .

Server:

http://localhost:8080

⸻

🐳 Docker

Build:

docker build -t order-api .

Run:

docker run -p 8080:8080 order-api

⸻

🗺️ Roadmap

flowchart LR
A[Go + Gin] --> B[Orders API]
B --> C[Payment System]
C --> D[PostgreSQL]
D --> E[Telegram Bot]
E --> F[Docker + CI/CD]

Phase 1 — Core API

* Gin server
* Order model
* In-memory storage
* Create / Read / Delete

Phase 2 — Payments

* POST /orders/:id/pay
* Payment service
* Stripe test integration
* Payment status
* Order state validation

Phase 3 — Database

* PostgreSQL
* Database schema
* Repository layer
* Migrations
* Replace in-memory storage

Phase 4 — Telegram

* Telegram bot
* Order creation
* Order lookup
* Payment button
* Payment notifications

Phase 5 — Infrastructure

* Dockerfile
* Docker Compose
* CI/CD
* Reverse proxy
* Deployment

⸻

🎯 Project Goal

Цель проекта — не просто сделать CRUD, а постепенно собрать полноценный backend:

        Telegram
           │
           ▼
      ┌─────────┐
      │ Go API  │
      │  Gin    │
      └────┬────┘
           │
     ┌─────┴─────┐
     ▼           ▼
PostgreSQL    Payment
Provider

Проект создаётся как практическая работа с Go, HTTP, бизнес-логикой, платежами, базами данных и инфраструктурой.
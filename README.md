OrderAPI

A small REST API for working with orders, written in Go.

This is a learning backend project where I am gradually building a proper application structure with configuration, HTTP server, PostgreSQL, Redis, Docker, and separate application layers.

Stack

* Go
* Gin
* PostgreSQL
* sqlx
* Redis
* Docker / Docker Compose
* YAML

Project structure

internal/
├── config/       # application configuration
├── domain/       # main models
├── handler/      # HTTP handlers
├── repository/   # database access
├── server/       # HTTP server and routes
└── service/      # business logic

Current state

* YAML configuration
* HTTP server with Gin
* GET /health
* Order domain model
* OrderRepository interface
* PostgreSQL repository with sqlx
* Creating orders in PostgreSQL
* Getting an order by ID

The project is still in development.

Running locally

Run the API:

make run

Run tests:

make test

Build the project:

make build

Format the code:

make form

Docker

Start the project:

make up

Stop the project:

make down

View logs:

make logs

Open PostgreSQL:

make db

TODO

* Finish the order repository
* Add the service layer
* Add CRUD endpoints
* Finish PostgreSQL integration
* Add database migrations
* Add Redis
* Add tests
* Finish Dockerfile and Docker Compose
* Add basic CI/CD

The project is being built step by step, so some parts are not implemented yet.

ifneq ("$(wildcard .env)","")
    include .env
    export $(shell sed 's/=.*//' .env)
endif

MIGRATIONS_DIR=./server/migrations
# Use DB_URL from .env as a fallback so you don't have to type it every time
URL ?= $(DB_URL)

## help: show this help message

run:
	@go run ./server/cmd/api

help:
	@echo Usage:
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/## //g' | column -t -s ':' || echo "Please install 'column' or just read the Makefile"

## migrate-add: add new migration files (Usage: make migrate-add name=create_users)
migrate-add:
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## migrate-up: apply all pending migrations
migrate-up:
	@if [ -z "$(URL)" ]; then echo "❌ Error: URL is not set. Use URL=... or set DB_URL in .env"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database $(URL) up

## migrate-down: revert one migration
migrate-down:
	@if [ -z "$(URL)" ]; then echo "❌ Error: URL is not set"; exit 1; fi
	migrate -path $(MIGRATIONS_DIR) -database $(URL) down 1

## migrate-status: show current migration version
migrate-status:
	migrate -path $(MIGRATIONS_DIR) -database $(URL) version

##Open a terminal inside the running database container and launch the PostgreSQL command-line client.
docker-repl:
	docker exec -it e2eechat-db psql -U e2ee -d e2eechat

.PHONY: help migrate-add migrate-up migrate-down migrate-status docker-repl
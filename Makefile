# Load .env so DB_* variables are available here and in every command below
-include .env
export

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATIONS_DIR := migrations

.PHONY: help run build test vet fmt tidy \
        db-up db-stop db-down db-logs db-psql \
        migrate-create migrate-up migrate-down migrate-version

help: ## Show available commands
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# ---------- App ----------
run: ## Run the API
	go run ./cmd/api

build: ## Build the binary into bin/api
	go build -o bin/api ./cmd/api

test: ## Run all tests
	go test ./...

vet: ## Static checks
	go vet ./...

fmt: ## Format all code
	go fmt ./...

tidy: ## Sync go.mod/go.sum with imports
	go mod tidy

# ---------- Database ----------
db-up: ## Start Postgres
	docker compose up -d

db-stop: ## Stop Postgres (keeps data)
	docker compose stop

db-down: ## Remove the Postgres container (keeps data volume)
	docker compose down

db-logs: ## Follow Postgres logs
	docker compose logs -f postgres

db-psql: ## Open a psql shell
	psql "$(DB_URL)"

# ---------- Migrations (needs the migrate CLI, step 2.1) ----------
migrate-create: ## New migration: make migrate-create name=create_users
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up: ## Apply all pending migrations
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down: ## Roll back the last migration
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-version: ## Show current migration version
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

dev: ## Run with hot reload (air)
	air --build.cmd "go build -o ./tmp/api ./cmd/api" --build.entrypoint "./tmp/api"
.PHONY: help run build clean test migrate-up migrate-down migrate-create migrate-force migrate-version seed-remove seed-reload docker-up docker-down swagger

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run cmd/server/main.go

build: ## Build the application
	go build -o bin/ewallet cmd/server/main.go

clean: ## Clean build artifacts
	rm -rf bin/

test: ## Run tests
	go test -v ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

swagger: ## Generate Swagger documentation
	~/go/bin/swag init -g cmd/server/main.go -o docs

migrate-up: ## Run database migrations
	~/go/bin/migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ewallet_db?sslmode=disable" up

migrate-down: ## Rollback database migrations
	~/go/bin/migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ewallet_db?sslmode=disable" down

migrate-create: ## Create new migration file (usage: make migrate-create name=migration_name)
	~/go/bin/migrate create -ext sql -dir migrations -seq $(name)

migrate-force: ## Force migration version (usage: make migrate-force version=1)
	~/go/bin/migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ewallet_db?sslmode=disable" force $(version)

migrate-version: ## Show current migration version
	~/go/bin/migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ewallet_db?sslmode=disable" version

seed-remove: ## Remove demo data (rollback seed migration)
	~/go/bin/migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/ewallet_db?sslmode=disable" down 1

seed-reload: ## Reload demo data (remove and apply again)
	@echo "Removing demo data..."
	@make seed-remove
	@echo "Reloading demo data..."
	@make migrate-up

docker-up: ## Start PostgreSQL with Docker Compose
	docker-compose up -d

docker-down: ## Stop Docker Compose services
	docker-compose down

docker-logs: ## View Docker Compose logs
	docker-compose logs -f

.DEFAULT_GOAL := help

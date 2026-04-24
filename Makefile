# Makefile for Workorder API

.PHONY: help build run test docker-up docker-down docker-logs clean

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o main .

run: ## Run the application locally
	go run main.go

test: ## Run tests
	go test -v ./...

docker-build: ## Build Docker image
	docker build -t workorder-api:latest .

docker-up: ## Start all services with Docker Compose
	docker-compose up -d --build

docker-up-dev: ## Start only dependencies (for local development)
	docker-compose -f docker-compose.dev.yml up -d

docker-down: ## Stop all services
	docker-compose down

docker-down-v: ## Stop all services and remove volumes (WARNING: deletes data!)
	docker-compose down -v

docker-logs: ## Show logs from all services
	docker-compose logs -f

docker-logs-api: ## Show logs from API service
	docker-compose logs -f api

docker-ps: ## Show running containers
	docker-compose ps

docker-restart: ## Restart all services
	docker-compose restart

docker-restart-api: ## Restart API service
	docker-compose restart api

clean: ## Clean build artifacts
	rm -f main
	go clean

deps: ## Download dependencies
	go mod download

tidy: ## Tidy dependencies
	go mod tidy

fmt: ## Format code
	go fmt ./...

lint: ## Run linter
	golangci-lint run

db-connect: ## Connect to PostgreSQL
	docker-compose exec postgres psql -U postgres -d workorder_db

db-backup: ## Backup database
	docker-compose exec postgres pg_dump -U postgres workorder_db > backup_$(shell date +%Y%m%d_%H%M%S).sql

rabbitmq-ui: ## Open RabbitMQ Management UI
	@echo "Opening RabbitMQ Management UI..."
	@echo "URL: http://localhost:15672"
	@echo "Username: guest"
	@echo "Password: guest"

api-health: ## Check API health
	curl http://localhost:8080/health

api-test: ## Test API with sample request
	curl -X POST http://localhost:8080/api/execute \
		-H "Content-Type: application/json" \
		-d @test_new_format.json

all: clean deps build ## Clean, download deps, and build

dev: docker-up-dev run ## Start dependencies and run API locally

prod: docker-up ## Start all services in production mode

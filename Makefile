.PHONY: help build up down restart logs clean test db-init db-reset ps health

# ========================================
# Default target
# ========================================
.DEFAULT_GOAL := help

# ========================================
# Variables
# ========================================
COMPOSE := docker-compose -f infrastructure/docker/docker-compose.yml
PROJECT_ROOT := $(shell pwd)
WEB_DIR := $(PROJECT_ROOT)/web/shop_mall_web

SERVICES_ALL := postgres redis rabbitmq mailhog \
	mock-stripe mock-fcm mock-elasticsearch mock-carriers \
	auth-service shop-service customer-service inventory-service \
	order-service payment-service shipping-service notification-service \
	review-service chat-service search-service admin-service

SERVICES_INFRA := postgres redis rabbitmq mailhog
SERVICES_MOCKS := mock-stripe mock-fcm mock-elasticsearch mock-carriers
SERVICES_MICRO := auth-service shop-service customer-service inventory-service \
	order-service payment-service shipping-service notification-service \
	review-service chat-service search-service admin-service
SERVICES_PHASE1 := auth-service shop-service

# ========================================
# Help
# ========================================
help: ## Show this help message
	@echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
	@echo '  Go Microservice Project - Makefile Commands'
	@echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'
	@echo ''
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@echo ''
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''
	@echo '━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'

# ========================================
# Build Commands
# ========================================
build: ## Build all services
	@echo "Building all services..."
	$(COMPOSE) build

build-infra: ## Build infrastructure services
	@echo "Building infrastructure services..."
	$(COMPOSE) build $(SERVICES_INFRA)

build-mocks: ## Build mock services
	@echo "Building mock services..."
	$(COMPOSE) build $(SERVICES_MOCKS)

build-services: ## Build microservices
	@echo "Building microservices..."
	$(COMPOSE) build $(SERVICES_MICRO)

build-no-cache: ## Build all services without cache
	@echo "Building all services without cache..."
	$(COMPOSE) build --no-cache

# ========================================
# Start/Stop Commands
# ========================================
up: ## Start all services
	@echo "Starting all services..."
	$(COMPOSE) up -d
	@echo "Waiting for services to be healthy..."
	@sleep 10
	@make ps

up-infra: ## Start infrastructure services only
	@echo "Starting infrastructure services..."
	$(COMPOSE) up -d $(SERVICES_INFRA)
	@make ps

up-mocks: ## Start mock services only
	@echo "Starting mock services..."
	$(COMPOSE) up -d $(SERVICES_MOCKS)
	@make ps

up-services: ## Start microservices only
	@echo "Starting microservices..."
	$(COMPOSE) up -d $(SERVICES_MICRO)
	@make ps

down: ## Stop all services
	@echo "Stopping all services..."
	$(COMPOSE) down

down-volumes: ## Stop all services and remove volumes
	@echo "Stopping all services and removing volumes..."
	$(COMPOSE) down -v

restart: ## Restart all services
	@echo "Restarting all services..."
	$(COMPOSE) restart

restart-infra: ## Restart infrastructure services
	@echo "Restarting infrastructure services..."
	$(COMPOSE) restart $(SERVICES_INFRA)

restart-services: ## Restart microservices
	@echo "Restarting microservices..."
	$(COMPOSE) restart $(SERVICES_MICRO)

# ========================================
# Status Commands
# ========================================
ps: ## Show running containers
	@$(COMPOSE) ps

logs: ## Show logs for all services
	$(COMPOSE) logs -f

logs-infra: ## Show logs for infrastructure services
	$(COMPOSE) logs -f $(SERVICES_INFRA)

logs-mocks: ## Show logs for mock services
	$(COMPOSE) logs -f $(SERVICES_MOCKS)

logs-services: ## Show logs for microservices
	$(COMPOSE) logs -f $(SERVICES_MICRO)

health: ## Check health status of all services
	@echo "Checking health status..."
	@$(COMPOSE) ps | grep -E "(healthy|running)" || true

# ========================================
# Phoenix Web Commands
# ========================================
phoenix: ## Start Phoenix web server (foreground)
	@echo "Starting Phoenix server..."
	@echo "Press Ctrl+C to stop"
	@cd $(WEB_DIR) && PORT=20200 mix phx.server

phoenix-bg: ## Start Phoenix web server (background)
	@echo "Starting Phoenix server in background..."
	@cd $(WEB_DIR) && PORT=20200 mix phx.server > /tmp/phoenix.log 2>&1 &
	@echo "✓ Phoenix started! (Logs: /tmp/phoenix.log)"
	@echo "  Access at: http://localhost:20200"

phoenix-stop: ## Stop Phoenix web server
	@echo "Stopping Phoenix server..."
	@pkill -f "mix phx.server" || echo "Phoenix not running"

phoenix-logs: ## Show Phoenix logs
	@tail -f /tmp/phoenix.log

phoenix-deps: ## Install Phoenix dependencies
	@echo "Installing Phoenix dependencies..."
	@cd $(WEB_DIR) && mix deps.get

phoenix-assets: ## Setup Phoenix assets
	@echo "Setting up Phoenix assets..."
	@cd $(WEB_DIR) && mix assets.setup

phoenix-compile: ## Compile Phoenix project
	@echo "Compiling Phoenix..."
	@cd $(WEB_DIR) && mix compile

phoenix-setup: ## Setup Phoenix (deps + assets)
	@echo "Setting up Phoenix project..."
	@make phoenix-deps
	@make phoenix-assets
	@make phoenix-compile
	@echo "✓ Phoenix setup complete!"

# ========================================
# Database Commands
# ========================================
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	@docker exec -i go_microservice_postgres_dev psql -U admin -d auth_db < scripts/migrations/auth/001_create_tables.sql
	@docker exec -i go_microservice_postgres_dev psql -U admin -d shop_db < scripts/migrations/shop/001_create_tables.sql
	@echo "✓ Migrations complete!"

db-init: ## Initialize databases
	@echo "Initializing databases..."
	@docker exec go_microservice_postgres_dev psql -U admin -d postgres -f /docker-entrypoint-initdb.d/init.sql

db-connect: ## Connect to PostgreSQL
	@docker exec -it go_microservice_postgres_dev psql -U admin -d postgres

db-list: ## List all databases
	@docker exec go_microservice_postgres_dev psql -U admin -d postgres -c "\l"

db-reset: ## Reset all databases (WARNING: removes all data)
	@echo "WARNING: This will remove all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		make down-volumes; \
		make up-infra; \
		sleep 15; \
		make db-init; \
		echo "Databases reset complete!"; \
	else \
		echo "Cancelled."; \
	fi

# ========================================
# Development Commands
# ========================================
dev: ## Start full development environment (infra + services + phoenix)
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "Starting full development environment..."
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@make up-infra
	@sleep 10
	@make db-migrate
	@make up-phase1
	@sleep 5
	@make phoenix-bg
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "✓ Development environment is ready!"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "Access URLs:"
	@echo "  - Phoenix Web:      http://localhost:20200/auth"
	@echo "  - Auth gRPC:        localhost:20100"
	@echo "  - Shop gRPC:        localhost:20101"
	@echo "  - PostgreSQL:       localhost:20000"
	@echo "  - Redis:            localhost:20001"
	@echo "  - RabbitMQ:         localhost:20002 (UI: http://localhost:20003)"
	@echo ""
	@echo "Useful commands:"
	@echo "  - View all logs:    make logs-all"
	@echo "  - View Phoenix:     make phoenix-logs"
	@echo "  - Stop all:         make dev-stop"
	@echo ""

dev-infra: ## Start development environment (infra + mocks only)
	@echo "Starting development infrastructure..."
	@make up-infra
	@sleep 15
	@make db-init
	@make up-mocks
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "Development infrastructure is ready!"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "Access URLs:"
	@echo "  - MailHog UI:       http://localhost:20005"
	@echo "  - RabbitMQ UI:      http://localhost:20003"
	@echo "  - PostgreSQL:       localhost:20000"
	@echo "  - Redis:            localhost:20001"
	@echo "  - Stripe Mock:      http://localhost:20010"
	@echo "  - FCM Mock:         http://localhost:20012"
	@echo "  - Elasticsearch:    http://localhost:20013"
	@echo "  - Carriers Mock:    http://localhost:20014"
	@echo ""

dev-stop: ## Stop development environment
	@echo "Stopping development environment..."
	@make down
	@make phoenix-stop
	@echo "✓ Development environment stopped!"

up-phase1: ## Start Phase 1 services (Auth + Shop)
	@echo "Starting Phase 1 services..."
	@$(COMPOSE) up -d $(SERVICES_PHASE1)
	@echo "✓ Phase 1 services started!"

logs-all: ## Show all logs (Docker + Phoenix)
	@echo "Showing all logs (Ctrl+C to stop)..."
	@( $(COMPOSE) logs -f & tail -f /tmp/phoenix.log 2>/dev/null & wait )

clean: ## Clean up containers, volumes, and images
	@echo "Cleaning up..."
	$(COMPOSE) down -v --remove-orphans
	@echo "Removing unused Docker images..."
	@docker image prune -f

test: ## Run tests (placeholder)
	@echo "Running tests..."
	@echo "TODO: Implement tests"

lint: ## Run linters (placeholder)
	@echo "Running linters..."
	@echo "TODO: Implement linters"

# ========================================
# Utility Commands
# ========================================
open-mailhog: ## Open MailHog UI in browser
	@open http://localhost:20005

open-rabbitmq: ## Open RabbitMQ UI in browser
	@open http://localhost:20003

shell-postgres: ## Open PostgreSQL shell
	@docker exec -it go_microservice_postgres_dev psql -U admin -d postgres

shell-redis: ## Open Redis CLI
	@docker exec -it go_microservice_redis_dev redis-cli -a redis_dev_password_123

backup-db: ## Backup all databases
	@echo "Backing up databases..."
	@mkdir -p backups
	@docker exec go_microservice_postgres_dev pg_dumpall -U admin > backups/backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "Backup complete: backups/backup_$$(date +%Y%m%d_%H%M%S).sql"

restore-db: ## Restore database from backup (usage: make restore-db FILE=backup_file.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "Error: Please specify FILE=backup_file.sql"; \
		exit 1; \
	fi
	@echo "Restoring database from $(FILE)..."
	@docker exec -i go_microservice_postgres_dev psql -U admin -d postgres < $(FILE)
	@echo "Restore complete!"

# ========================================
# Code Generation Commands
# ========================================
generate: ## Generate code from DSL
	@echo "Generating code from DSL..."
	@./scripts/mps-generate.sh --all

generate-proto: ## Generate Protocol Buffers
	@echo "Generating Protocol Buffers..."
	@echo "TODO: Implement proto generation"

go-tidy: ## Run go mod tidy for all services
	@echo "Running go mod tidy for all services..."
	@for dir in generated/*/; do \
		echo "Processing $$dir"; \
		cd "$$dir" && go mod tidy && cd - > /dev/null; \
	done
	@for dir in mock/*/; do \
		echo "Processing $$dir"; \
		cd "$$dir" && go mod tidy && cd - > /dev/null; \
	done
	@echo "Done!"

# ========================================
# Monitoring Commands
# ========================================
stats: ## Show container resource usage
	@docker stats --no-stream $$(docker ps --format "{{.Names}}" | grep go_microservice)

top: ## Show running processes in containers
	@echo "Showing top processes in containers..."
	@for container in $$(docker ps --format "{{.Names}}" | grep go_microservice | head -5); do \
		echo ""; \
		echo "=== $$container ==="; \
		docker top $$container; \
	done

inspect: ## Inspect a service (usage: make inspect SERVICE=postgres)
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: Please specify SERVICE=service_name"; \
		exit 1; \
	fi
	@$(COMPOSE) exec $(SERVICE) env

# ========================================
# Quick Start
# ========================================
init: ## Initial setup (first time only)
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "Go MicroService Example - Initial Setup"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@if [ ! -f .env ]; then \
		echo "Creating .env file..."; \
		cp .env.example .env; \
		echo "✓ .env created"; \
	else \
		echo "✓ .env already exists"; \
	fi
	@echo ""
	@echo "Setting up Phoenix..."
	@make phoenix-setup
	@echo ""
	@echo "Building Docker services..."
	@make build-mocks
	@$(COMPOSE) build $(SERVICES_PHASE1)
	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "✓ Initial setup complete!"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Start all services:  make dev"
	@echo "  2. Open browser:        open http://localhost:20200/auth"
	@echo "  3. Check status:        make status"
	@echo "  4. View logs:           make logs-all"
	@echo ""

quickstart: ## Quick start (setup + run)
	@make init
	@make dev

status: ## Show system status
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "System Status"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "Docker Containers:"
	@$(COMPOSE) ps --format "table {{.Name}}\t{{.Status}}" | head -10
	@echo ""
	@echo -n "Phoenix Server: "
	@pgrep -f "mix phx.server" > /dev/null 2>&1 && echo "✓ Running" || echo "✗ Not running"
	@echo ""
	@echo "Access URLs:"
	@echo "  - Phoenix Web:      http://localhost:20200/auth"
	@echo "  - Auth gRPC:        localhost:20100"
	@echo "  - Shop gRPC:        localhost:20101"
	@echo ""

# ========================================
# Danger Zone
# ========================================
nuke: ## Remove everything (containers, volumes, images) - USE WITH CAUTION
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "⚠️  WARNING: This will remove EVERYTHING!"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo ""
	@echo "This will delete:"
	@echo "  - All containers"
	@echo "  - All volumes (ALL DATA WILL BE LOST)"
	@echo "  - All project Docker images"
	@echo ""
	@read -p "Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "Removing everything..."; \
		$(COMPOSE) down -v --remove-orphans; \
		docker images | grep go_microservice | awk '{print $$3}' | xargs -r docker rmi -f; \
		echo "✓ Everything removed!"; \
	else \
		echo "Cancelled."; \
	fi

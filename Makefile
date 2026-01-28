.PHONY: help up down status logs clean test build

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "🚀 Go Microservice - Available Commands"
	@echo "========================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

up: ## Start all services (databases + microservices)
	@echo "🚀 Starting all services..."
	@cd infrastructure/docker && docker compose up -d
	@sleep 5
	@echo "✅ Databases started"
	@make start-services
	@echo ""
	@make status

down: ## Stop all services
	@echo "🛑 Stopping all services..."
	@make stop-services
	@cd infrastructure/docker && docker compose down
	@echo "✅ All services stopped"

start-services: ## Start all microservices
	@echo "🔧 Starting microservices..."
	@./scripts/start_all_services.sh

stop-services: ## Stop all microservices
	@echo "🛑 Stopping microservices..."
	@./scripts/stop_all_services.sh

restart: down up ## Restart all services

status: ## Show status of all services
	@echo "📊 Service Status"
	@echo "================="
	@./scripts/check_all_services.sh

logs: ## Show logs for all services
	@echo "📋 Service Logs"
	@echo "==============="
	@tail -20 /tmp/*-service.log 2>/dev/null || echo "No logs found"
	@tail -20 microservices/auth/auth-server.log 2>/dev/null || echo "No auth logs"

logs-follow: ## Follow logs in real-time
	@tail -f /tmp/*-service.log microservices/auth/auth-server.log 2>/dev/null

clean: ## Clean up logs and temporary files
	@echo "🧹 Cleaning up..."
	@rm -f /tmp/*-service.log
	@rm -f microservices/auth/auth-server.log
	@echo "✅ Cleanup complete"

test: ## Run all integration tests
	@echo "🧪 Running all tests..."
	@cd tests && ./run_all_integration_tests.sh

test-auth: ## Run auth integration tests
	@echo "🔐 Running auth tests..."
	@cd tests/integration/auth && ./run_test.sh

test-order: ## Run order flow tests
	@echo "📦 Running order flow tests..."
	@cd tests/integration/order_flow && ./run_test.sh

test-e2e: ## Run E2E tests
	@echo "🎯 Running E2E tests..."
	@cd tests/e2e && ./test_runner.sh

build: ## Build all services
	@echo "🔨 Building all services..."
	@./scripts/build_all_services.sh

db-status: ## Show database status
	@echo "🗄️  Database Status"
	@echo "=================="
	@docker ps --filter "name=postgres" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

db-logs: ## Show database logs
	@docker compose -f infrastructure/docker/docker-compose.yml logs postgres_auth postgres_shop postgres_customer

ps: ## Show running processes
	@echo "🔍 Running Processes"
	@echo "===================="
	@ps aux | grep -E "auth-server|shop-server|customer-service|inventory-service|order-server|payment-server|notification-service|review-service|shipping|chat-service|search-service|admin-service" | grep -v grep | awk '{print $$2, $$11}' || echo "No services running"

dashboard: ## Open service dashboard
	@echo "📊 Service Dashboard"
	@cat RUNNING_SERVICES_DASHBOARD.md

dev: ## Start in development mode (with logs)
	@make up
	@make logs-follow

.PHONY: help build run clean test db-test

help:
	@echo "Search Service - Makefile commands"
	@echo ""
	@echo "Available commands:"
	@echo "  make build     - Build the service binary"
	@echo "  make run       - Run the service"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make test      - Run tests"
	@echo "  make db-test   - Test database connection"

build:
	@echo "Building Search Service..."
	go build -o search-service main.go
	@echo "✅ Build complete: search-service"

run:
	@echo "Starting Search Service..."
	go run main.go

clean:
	@echo "Cleaning..."
	rm -f search-service
	go clean
	@echo "✅ Clean complete"

test:
	@echo "Running tests..."
	go test -v ./...

db-test:
	@echo "Testing database connection..."
	@PGPASSWORD=postgres_password psql -h localhost -p 22020 -U postgres -d search_service -c "SELECT 'Connection successful!' as status;"

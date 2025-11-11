.PHONY: run dev tidy test test-order test-coverage swagger mock \
        migrate-up migrate-down migrate-status migrate-create db-setup db-reset \
        docker-up docker-down docker-logs docker-restart \
        prepare

SWAG_BIN := $(shell go env GOPATH)/bin/swag
MOCKGEN_BIN := $(shell go env GOPATH)/bin/mockgen

prepare:
	@echo "Installing development tools..."
	@echo "→ Installing Air (hot reload)..."
	@go install github.com/air-verse/air@latest

	@echo "→ Installing Swag (swagger)..."
	@go install github.com/swaggo/swag/cmd/swag@latest

	@echo "→ Installing Mockgen..."
	@go install go.uber.org/mock/mockgen@latest

	@echo "✓ All tools installed successfully!"

tidy:
	@echo "Tidying go modules..."
	@go mod tidy
	@echo "✓ Go modules tidied!"

run:
	@echo "Starting application..."
	@go run .

dev:
	@echo "Starting application with hot reload..."
	@air || go run .

test:
	@echo "Running tests..."
	@go test -v -race -cover ./feature/order/...

test-coverage:
	@echo "Running tests with coverage report..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

db-setup: migrate-up
	@echo "✓ Database setup completed!"

db-reset:
	@echo "Resetting database..."
	@go run scripts/migrate/main.go -cmd=down || true
	@go run scripts/migrate/main.go -cmd=down || true
	@go run scripts/migrate/main.go -cmd=down || true
	@$(MAKE) migrate-up
	@echo "✓ Database reset completed!"

migrate-up:
	@echo "Running migrations..."
	@go run scripts/migrate/main.go -cmd=up

migrate-down:
	@echo "Rolling back migration..."
	@go run scripts/migrate/main.go -cmd=down

migrate-status:
	@echo "Checking migration status..."
	@go run scripts/migrate/main.go -cmd=status

migrate-create:
	@echo "Creating new migration..."
	@go run scripts/generate_migration/main.go

docker-up:
	@echo "Starting docker services..."
	@docker compose up -d
	@echo "✓ Docker services started!"
	@echo "→ PostgreSQL: localhost:5432"

docker-down:
	@echo "Stopping docker services..."
	@docker compose down
	@echo "✓ Docker services stopped!"

docker-logs:
	@echo "Showing docker logs..."
	@docker compose logs -f

docker-restart:
	@echo "Restarting docker services..."
	@docker compose restart
	@echo "✓ Docker services restarted!"

swagger:
	@echo "Generating swagger docs..."
	@$(SWAG_BIN) fmt
	@$(SWAG_BIN) init -g main.go -o docs
	@echo "✓ Swagger docs generated!"

mock:
	@echo "Generating mocks..."
	@$(MOCKGEN_BIN) -source=domain/order.go -destination=domain/mock/mock_order.go -package=mock
	@$(MOCKGEN_BIN) -source=domain/user.go -destination=domain/mock/mock_user.go -package=mock
	@$(MOCKGEN_BIN) -source=domain/cart.go -destination=domain/mock/mock_cart.go -package=mock
	@$(MOCKGEN_BIN) -source=domain/product.go -destination=domain/mock/mock_product.go -package=mock
	@$(MOCKGEN_BIN) -source=domain/shop.go -destination=domain/mock/mock_shop.go -package=mock
	@echo "✓ Mocks generated successfully!"

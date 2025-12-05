.PHONY: help build run test clean docker-up docker-down db-create db-drop migrate swagger lint

# Переменные
APP_NAME=wallet_server
GO_FILES=$(shell find . -name "*.go" -type f)
MAIN_FILE=main.go

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Показать эту справку
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

build: ## Собрать приложение
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	go build -o $(APP_NAME) $(MAIN_FILE)
	@echo "$(GREEN)✓ Build complete: ./$(APP_NAME)$(NC)"

run: ## Запустить приложение
	@echo "$(GREEN)Starting $(APP_NAME)...$(NC)"
	go run $(MAIN_FILE)

dev: ## Запустить в режиме разработки (с автоперезагрузкой)
	@echo "$(GREEN)Starting in dev mode...$(NC)"
	@which air > /dev/null || (echo "$(RED)air not installed. Run: go install github.com/cosmtrek/air@latest$(NC)" && exit 1)
	air

test: ## Запустить тесты
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

clean: ## Очистить сборочные файлы
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html
	rm -rf docs/swagger/
	@echo "$(GREEN)✓ Cleaned$(NC)"

install: ## Установить зависимости
	@echo "$(GREEN)Installing dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

lint: ## Запустить линтер
	@echo "$(GREEN)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed$(NC)" && exit 1)
	golangci-lint run ./...

fmt: ## Форматировать код
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

# База данных
db-create: ## Создать базу данных
	@echo "$(GREEN)Creating database...$(NC)"
	createdb wallet_db || echo "$(YELLOW)Database might already exist$(NC)"
	@echo "$(GREEN)✓ Database created$(NC)"

db-drop: ## Удалить базу данных
	@echo "$(RED)Dropping database...$(NC)"
	dropdb wallet_db || echo "$(YELLOW)Database might not exist$(NC)"

db-reset: db-drop db-create ## Пересоздать базу данных
	@echo "$(GREEN)✓ Database reset$(NC)"

migrate: ## Запустить миграции (выполняется автоматически при старте)
	@echo "$(GREEN)Migrations run automatically on app start$(NC)"

# Docker
docker-build: ## Собрать Docker образ
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t wallet_test:latest .
	@echo "$(GREEN)✓ Docker image built$(NC)"

docker-up: ## Запустить через Docker Compose
	@echo "$(GREEN)Starting Docker containers...$(NC)"
	docker compose up -d
	@echo "$(GREEN)✓ Containers started$(NC)"

docker-down: ## Остановить Docker контейнеры
	@echo "$(YELLOW)Stopping Docker containers...$(NC)"
	docker compose down
	@echo "$(GREEN)✓ Containers stopped$(NC)"

docker-logs: ## Показать логи Docker контейнеров
	docker compose logs -f

# Swagger
swagger: ## Сгенерировать Swagger документацию
	@echo "$(GREEN)Generating Swagger docs...$(NC)"
	@which swag > /dev/null || (echo "$(RED)swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest$(NC)" && exit 1)
	swag init -g src/main.go --output docs/swagger
	@echo "$(GREEN)✓ Swagger docs generated in docs/swagger/$(NC)"

# Тестирование API
test-api: ## Тестировать API endpoints
	@echo "$(GREEN)Testing API endpoints...$(NC)"
	bash examples.sh

# Проверки
check: fmt lint test ## Выполнить все проверки

# Установка инструментов разработки
tools: ## Установить инструменты разработки
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Tools installed$(NC)"

# Информация
info: ## Показать информацию о проекте
	@echo "$(GREEN)Project: TON Wallet API$(NC)"
	@echo "  Go version:    $$(go version | awk '{print $$3}')"
	@echo "  Binary:        $(APP_NAME)"
	@echo "  Main file:     $(MAIN_FILE)"
	@echo "  Go files:      $$(echo $(GO_FILES) | wc -w)"

# Запуск примеров
example-wallet: ## Запустить пример работы с TON кошельком
	@echo "$(GREEN)Running wallet example...$(NC)"
	cd wallet && go run main.go

# Установка окружения
setup: install db-create ## Полная настройка окружения
	@echo "$(GREEN)Setting up environment...$(NC)"
	@test -f .env || (echo "Creating .env from .env.example..." && cp .env.example .env)
	@echo "$(GREEN)✓ Setup complete!$(NC)"
	@echo ""
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Edit .env file with your settings"
	@echo "  2. Run: make build"
	@echo "  3. Run: make run"

.DEFAULT_GOAL := help

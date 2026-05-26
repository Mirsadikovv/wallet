.PHONY: build run dev test clean install lint fmt swagger check \
        db-create db-drop db-reset migrate

APP_NAME=wallet_server
ROOT=..

GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m

# ─── Сборка ────────────────────────────────────────────────────────────────────

build: ## Собрать бинарник
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	go build -o ../$(APP_NAME) .
	@echo "$(GREEN)✓ $(ROOT)/$(APP_NAME)$(NC)"

run: ## Запустить
	@echo "$(GREEN)Starting $(APP_NAME)...$(NC)"
	go run .

dev: ## Запустить с hot-reload (air)
	@echo "$(GREEN)Starting dev mode...$(NC)"
	@which air > /dev/null || (echo "$(RED)air not installed: go install github.com/cosmtrek/air@latest$(NC)" && exit 1)
	air

# ─── Качество ──────────────────────────────────────────────────────────────────

test: ## Тесты с race detector и покрытием
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v -race -coverprofile=$(ROOT)/coverage.out ./...
	go tool cover -html=$(ROOT)/coverage.out -o $(ROOT)/coverage.html
	@echo "$(GREEN)✓ coverage.html$(NC)"

lint: ## Линтер
	@echo "$(GREEN)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed$(NC)" && exit 1)
	golangci-lint run ./...

fmt: ## Форматирование
	@echo "$(GREEN)Formatting...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✓ Done$(NC)"

check: fmt lint test ## fmt + lint + test

# ─── Зависимости ───────────────────────────────────────────────────────────────

install: ## Скачать и привести в порядок зависимости
	@echo "$(GREEN)Installing dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(GREEN)✓ Done$(NC)"

# ─── База данных ───────────────────────────────────────────────────────────────

db-create: ## Создать базу данных
	@echo "$(GREEN)Creating database...$(NC)"
	createdb wallet_db || echo "$(YELLOW)Already exists$(NC)"

db-drop: ## Удалить базу данных
	@echo "$(RED)Dropping database...$(NC)"
	dropdb wallet_db || echo "$(YELLOW)Not found$(NC)"

db-reset: db-drop db-create ## Пересоздать базу данных

migrate: ## Миграции запускаются автоматически при старте
	@echo "$(GREEN)Migrations run automatically on app start$(NC)"

# ─── Документация ──────────────────────────────────────────────────────────────

swagger: ## Сгенерировать Swagger
	@echo "$(GREEN)Generating Swagger...$(NC)"
	@which swag > /dev/null 2>&1 || test -f $(shell go env GOPATH)/bin/swag || \
		(echo "$(RED)swag not installed: go install github.com/swaggo/swag/cmd/swag@latest$(NC)" && exit 1)
	$(shell which swag 2>/dev/null || echo $(shell go env GOPATH)/bin/swag) init -g src/main.go --output ./docs/swagger
	@echo "$(GREEN)✓ backend/docs/swagger/$(NC)"

# ─── Утилиты ───────────────────────────────────────────────────────────────────

clean: ## Удалить артефакты сборки
	rm -f $(ROOT)/$(APP_NAME) $(ROOT)/coverage.out $(ROOT)/coverage.html
	@echo "$(GREEN)✓ Cleaned$(NC)"

.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-18s\033[0m %s\n", $$1, $$2}'

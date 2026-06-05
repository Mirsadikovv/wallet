.PHONY: help \
        be.build be.run be.dev be.test be.lint be.fmt be.check be.swagger be.install be.clean \
        be.db-create be.db-drop be.db-reset \
        fe.install fe.dev fe.build fe.lint fe.fmt fe.check fe.clean \
        bot.run bot.dev bot.build bot.install bot.clean \
        docker-build docker-up docker-down docker-restart docker-logs docker-ps \
        tools setup info test-api clean

APP_NAME=wallet_server

GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m

# ═══════════════════════════════════════════════════════════════════════════════
#  BACKEND  (be.*)
# ═══════════════════════════════════════════════════════════════════════════════

be.build: ## Backend: собрать бинарник
	$(MAKE) -C backend build

be.run: ## Backend: запустить
	$(MAKE) -C backend run

be.dev: ## Backend: hot-reload (air)
	$(MAKE) -C backend dev

be.test: ## Backend: тесты
	$(MAKE) -C backend test

be.lint: ## Backend: линтер
	$(MAKE) -C backend lint

be.fmt: ## Backend: форматирование
	$(MAKE) -C backend fmt

be.check: ## Backend: fmt + lint + test
	$(MAKE) -C backend check

be.swagger: ## Backend: сгенерировать Swagger
	$(MAKE) -C backend swagger

be.install: ## Backend: установить зависимости
	$(MAKE) -C backend install

be.clean: ## Backend: удалить артефакты
	$(MAKE) -C backend clean

be.db-create: ## Backend: создать БД
	$(MAKE) -C backend db-create

be.db-drop: ## Backend: удалить БД
	$(MAKE) -C backend db-drop

be.db-reset: ## Backend: пересоздать БД
	$(MAKE) -C backend db-reset

# ═══════════════════════════════════════════════════════════════════════════════
#  FRONTEND  (fe.*)
# ═══════════════════════════════════════════════════════════════════════════════

fe.install: ## Frontend: установить зависимости
	$(MAKE) -C frontend install

fe.dev: ## Frontend: dev-сервер
	$(MAKE) -C frontend dev

fe.build: ## Frontend: продакшен-сборка
	$(MAKE) -C frontend build

fe.lint: ## Frontend: линтер
	$(MAKE) -C frontend lint

fe.fmt: ## Frontend: форматирование
	$(MAKE) -C frontend fmt

fe.check: ## Frontend: все проверки
	$(MAKE) -C frontend check

fe.clean: ## Frontend: удалить артефакты
	$(MAKE) -C frontend clean

# ═══════════════════════════════════════════════════════════════════════════════
#  BOT  (bot.*)
# ═══════════════════════════════════════════════════════════════════════════════

bot.install: ## Bot: установить зависимости
	cd bot && go mod tidy

bot.run: ## Bot: запустить (prod)
	cd bot && go run .

bot.dev: ## Bot: запустить с hot-reload (dev)
	cd bot && go run -tags dev .

bot.build: ## Bot: собрать бинарник
	cd bot && go build -tags dev -o ../wallet_bot .

bot.clean: ## Bot: удалить бинарник
	rm -f wallet_bot

# ═══════════════════════════════════════════════════════════════════════════════
#  DOCKER
# ═══════════════════════════════════════════════════════════════════════════════

docker-build: ## Docker: собрать все образы
	@echo "$(GREEN)Building all Docker images...$(NC)"
	docker compose build
	@echo "$(GREEN)✓ Done$(NC)"

docker-up: ## Docker: запустить все контейнеры
	@echo "$(GREEN)Starting containers...$(NC)"
	docker compose up -d
	@echo "$(GREEN)✓ Done$(NC)"

docker-down: ## Docker: остановить контейнеры
	@echo "$(YELLOW)Stopping containers...$(NC)"
	docker compose down
	@echo "$(GREEN)✓ Done$(NC)"

docker-restart: ## Docker: пересобрать и перезапустить
	@echo "$(YELLOW)Rebuilding and restarting...$(NC)"
	docker compose down
	docker compose build
	docker compose up -d
	@echo "$(GREEN)✓ Done$(NC)"

docker-logs: ## Docker: логи всех сервисов
	docker compose logs -f

docker-ps: ## Docker: статус контейнеров
	docker compose ps

# ═══════════════════════════════════════════════════════════════════════════════
#  ОБЩИЕ
# ═══════════════════════════════════════════════════════════════════════════════

clean: be.clean fe.clean ## Очистить все артефакты
	rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ All cleaned$(NC)"

test-api: ## Тестировать API через examples.sh
	@echo "$(GREEN)Testing API...$(NC)"
	bash examples.sh

tools: ## Установить инструменты разработки
	@echo "$(GREEN)Installing dev tools...$(NC)"
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Done$(NC)"

setup: be.install be.db-create ## Полная настройка окружения
	@echo "$(GREEN)Setting up environment...$(NC)"
	@test -f .env || (echo "$(YELLOW)Creating .env from .env.example...$(NC)" && cp .env.example .env)
	@echo "$(GREEN)✓ Setup complete$(NC)"
	@echo ""
	@echo "$(YELLOW)Next steps:$(NC)"
	@echo "  1. Edit .env"
	@echo "  2. make be.run"

info: ## Информация о проекте
	@echo "$(GREEN)TON Wallet$(NC)"
	@echo "  Go:       $$(go version | awk '{print $$3}')"
	@echo "  Binary:   ./$(APP_NAME)"
	@echo "  Backend:  backend/"
	@echo "  Frontend: frontend/"

# ═══════════════════════════════════════════════════════════════════════════════
#  АЛИАСЫ (обратная совместимость)
# ═══════════════════════════════════════════════════════════════════════════════

build:   be.build
run:     be.run
dev:     be.dev
test:    be.test
lint:    be.lint
fmt:     be.fmt
check:   be.check
swagger: be.swagger
install: be.install
db-create: be.db-create
db-drop:   be.db-drop
db-reset:  be.db-reset

# ═══════════════════════════════════════════════════════════════════════════════

.DEFAULT_GOAL := help

help:
	@echo ""
	@echo "$(GREEN)╔══════════════════════════════════════╗$(NC)"
	@echo "$(GREEN)║         TON Wallet — Makefile        ║$(NC)"
	@echo "$(GREEN)╚══════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(YELLOW)Backend (be.*):$(NC)"
	@grep -E '^be\.[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Frontend (fe.*):$(NC)"
	@grep -E '^fe\.[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Bot (bot.*):$(NC)"
	@grep -E '^bot\.[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Docker:$(NC)"
	@grep -E '^docker-[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Общие:$(NC)"
	@grep -E '^(clean|test-api|tools|setup|info):.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'
	@echo ""

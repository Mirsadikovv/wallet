.PHONY: install dev build lint fmt check clean

GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m

PKG_MANAGER=bun

# ─── Зависимости ───────────────────────────────────────────────────────────────

install: ## Установить зависимости
	@echo "$(GREEN)Installing frontend dependencies...$(NC)"
	$(PKG_MANAGER) install
	@echo "$(GREEN)✓ Done$(NC)"

# ─── Разработка ────────────────────────────────────────────────────────────────

dev: ## Запустить dev-сервер (Vite)
	@echo "$(GREEN)Starting frontend dev server...$(NC)"
	$(PKG_MANAGER) run dev

# ─── Сборка ────────────────────────────────────────────────────────────────────

build: ## Собрать для продакшена
	@echo "$(GREEN)Building frontend...$(NC)"
	$(PKG_MANAGER) run build
	@echo "$(GREEN)✓ dist/$(NC)"

# ─── Качество ──────────────────────────────────────────────────────────────────

lint: ## Запустить линтер (ESLint)
	@echo "$(GREEN)Running ESLint...$(NC)"
	$(PKG_MANAGER) run lint

fmt: ## Форматировать (Prettier)
	@echo "$(GREEN)Formatting...$(NC)"
	$(PKG_MANAGER) run fmt

check: lint ## Все проверки

# ─── Утилиты ───────────────────────────────────────────────────────────────────

clean: ## Удалить артефакты сборки
	rm -rf dist node_modules/.cache
	@echo "$(GREEN)✓ Cleaned$(NC)"

.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-18s\033[0m %s\n", $$1, $$2}'

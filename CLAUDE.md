# Wallet — Claude Context

## Стек
- **Backend**: Go 1.25+, Echo v4, GORM, PostgreSQL, Redis
- **Frontend**: Vue 3, TypeScript, Quasar 2, Pinia, UnoCSS
- **Инфра**: Docker Compose, Makefile

## Структура репо
```
wallet/
├── backend/          ← Go backend
│   └── src/
│       ├── modules/  ← доменные модули
│       └── common/   ← enum, helpers, middleware, utils
├── frontend/         ← Vue 3 frontend
└── Makefile          ← все команды через make
```

## Частые команды
```
make dev          ← запуск backend с hot-reload (air)
make run          ← запуск без hot-reload
make test         ← тесты
make swagger      ← регенерировать Swagger документацию
make fmt          ← gofmt
make lint         ← линтер
make check        ← fmt + lint + test
make docker-up    ← запустить через Docker Compose
make migrate      ← запустить миграции
```

## Стайл-гайды (читай перед написанием кода)
@.claude/CLAUDE_STYLE_BACKEND.md
@.claude/CLAUDE_STYLE_FRONTEND.md
@.claude/CLAUDE_STYLE_DEVOPS.md

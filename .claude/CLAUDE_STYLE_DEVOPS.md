# CLAUDE_STYLE — DevOps + Overall

> Структура репозитория, инфраструктура, окружения.

---

## Структура репозитория

```
/idoctor_platform/                  ← корень
│
├── idoctor_back/                   ← Go-бэкенд (API)
│   ├── src/
│   │   ├── module/                 ← модули по доменам
│   │   ├── common/                 ← общий код (enum, helpers, middleware, utils...)
│   │   └── main.go
│   ├── tests/                      ← тесты (на одном уровне с src/)
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── dev.go
│   └── prod.go
│
├── idoctor_front/                  ← Vue 3 фронтенд
│   ├── src/
│   ├── public/
│   ├── package.json
│   ├── vite.config.ts
│   └── Dockerfile
│
├── idoctor_bot/                    ← Telegram-бот (Go)
│   ├── app/
│   ├── main.go
│   ├── dev.go
│   └── prod.go
│
├── docker-compose.yml              ← продакшн
├── docker-compose-dev.yml          ← разработка
├── docker-compose-prod.yml         ← продакшн (алтернативный)
├── Dockerfile                      ← для бэкенда
├── Makefile
└── .github/
    └── workflows/
        ├── dev.yml
        └── prod.yml
```

**Правило**: Docker и CI/CD файлы живут **только в корне**. Внутри `idoctor_back/`, `idoctor_front/`, `idoctor_bot/` — только код.

---

## Docker Compose сервисы

| Сервис | Образ | Назначение |
|--------|-------|------------|
| `traefik` | traefik:v2.11 | Reverse proxy + SSL (Let's Encrypt) |
| `db-platform-idoctor` | postgres:16-alpine | Основная БД |
| `platform-idoctor-redis-db` | redis:7.2.4-alpine | Кеш |
| `platform-idoctor-app` | (build) | Go API |
| `platform_idoctor_bot` | (build) | Telegram бот |
| `idoctor_dev_dashboard` | (build) | Vue фронтенд |

---

## Окружения

| Файл | Назначение |
|------|------------|
| `docker-compose.yml` | Базовая конфигурация |
| `docker-compose-dev.yml` | Dev-окружение |
| `docker-compose-prod.yml` | Prod-окружение |
| `.github/workflows/dev.yml` | CI/CD для dev ветки |
| `.github/workflows/prod.yml` | CI/CD для prod ветки |

---

## Переменные окружения (бэкенд)

Все env-переменные описываются в Go-структуре с тегами `env`:

```go
type Env struct {
    HTTP_Host          string `env:"HTTP_HOST" default:"localhost"`
    HTTP_Port          int    `env:"HTTP_PORT" default:"80"`
    POSTGRES_Host      string `env:"POSTGRES_HOST"`
    POSTGRES_User      string `env:"POSTGRES_USER"`
    POSTGRES_Password  string `env:"POSTGRES_PASSWORD"`
    POSTGRES_DBName    string `env:"POSTGRES_DB"`
    POSTGRES_Port      int    `env:"POSTGRES_PORT"`
    POSTGRES_SSLMode   string `env:"POSTGRES_SSL_MODE" default:"disable"`
    POSTGRES_TimeZone  string `env:"POSTGRES_TIME_ZONE" default:"UTC"`
    JWT_Secret         string `env:"JWT_SECRET"`
    JWT_Expired        int64  `env:"JWT_EXPIRED"`
    JWT_RefreshExpired int64  `env:"JWT_REFRESH_EXPIRED"`
    REDIS_Addr         string `env:"REDIS_ADDR"`
}
```

Переменные окружения (фронтенд) — через `.env`:

```
VITE_API_URL=...
VITE_ONEID_*=...
VITE_EIMZO_*=...
```

---

## Makefile

Команды запускаются через `make`. Не вызывать `go run`, `docker compose` напрямую — использовать таргеты из `Makefile`.

---

## CI/CD

- **dev ветка** → `dev.yml` → деплой на dev-окружение
- **prod ветка** → `prod.yml` → деплой на prod-окружение
- Пайплайн: сборка Docker-образа → push → деплой на сервер

---

## Правила для корневых файлов

- `Dockerfile` в корне — для **бэкенда** (`idoctor_back/`)
- `idoctor_front/Dockerfile` — для **фронтенда**
- `idoctor_bot/` не имеет отдельного Dockerfile (запускается через общий compose)
- `.github/` — только CI/CD workflows, ничего лишнего

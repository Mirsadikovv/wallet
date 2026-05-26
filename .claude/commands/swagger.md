# swagger

Для выделенного хендлера напиши полные Swagger godoc аннотации.

## Правила

- Используй шаблон из CLAUDE_STYLE_BACKEND.md
- Теги: `@Summary`, `@Description`, `@Tags`, `@Accept`, `@Produce`, `@Security`, `@Param`, `@Success`, `@Failure`, `@Router`
- `@Security ApiKeyAuth` — для всех защищённых роутов
- `@Success 200` или `201` в зависимости от метода
- `@Failure 400` и `@Failure 500` обязательны
- Не меняй тело функции — только добавляй комментарии выше `func`

## HTTP коды по методам
- POST (create) → 201
- GET → 200
- PUT / PATCH → 200
- DELETE / restore → 200

# new-module

Создай полный Go-модуль с именем `$ARGUMENTS` в `backend/src/modules/`.

## Что создать

Структура папок и файлов:
```
backend/src/modules/$ARGUMENTS_service/
├── cmd.go
├── handler/$ARGUMENTS_handler.go
├── service/$ARGUMENTS_service.go
├── dto/$ARGUMENTS_dto.go
└── model/$ARGUMENTS_model.go
```

## Требования

- Строго следуй правилам из CLAUDE_STYLE_BACKEND.md
- CRUD методы: `Page`, `FindByID`, `Create`, `Update`, `DeleteOrRestore`
- Интерфейс + приватная реализация (`NewXxxService` возвращает интерфейс)
- Swagger godoc на каждом хендлере (обязателен)
- DTO: теги `json` + `validate`, аннотация `// @Name`
- Model: минимальный GORM с `DeletedAt` для soft delete
- Анонимные блоки `{}` для обработки ошибок в хендлерах
- Никаких комментариев кроме Swagger godoc

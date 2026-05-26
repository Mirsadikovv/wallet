# CLAUDE_STYLE — Backend (Go)

> Правила для кода в `idoctor_back/`. Строго следуй при написании нового кода.

---

## Стек

- **Язык**: Go 1.25+
- **Web**: Echo v4
- **ORM**: GORM + PostgreSQL
- **Cache**: Redis
- **JWT**: golang-jwt/jwt/v5
- **Shared библиотека**: `github.com/Mirsadikovv/shared` — `pg`, `logger`, `request`, `response`, `jwt`, `redis`
- **Документация**: Swagger (godoc аннотации)

---

## Структура

```
idoctor_back/
├── src/
│   ├── module/
│   │   └── <name>_service/
│   │       ├── cmd.go                    ← регистрация роутов
│   │       ├── handler/<name>_handler.go ← HTTP обработчики
│   │       ├── service/<name>_service.go ← бизнес-логика + интерфейс
│   │       ├── dto/<name>_dto.go         ← структуры запрос/ответ
│   │       └── model/<name>_model.go     ← GORM модель
│   ├── common/
│   │   ├── enum/        ← типизированные enum'ы
│   │   ├── helpers/     ← вспомогательные функции
│   │   ├── middleware/  ← Echo middleware
│   │   ├── seeder/      ← сидеры БД
│   │   ├── static/      ← константы и сообщения
│   │   ├── utils/       ← утилиты
│   │   └── validator/   ← кастомные валидаторы
│   └── main.go
├── tests/               ← тесты (на одном уровне с src/)
│   └── <name>_service_test.go
├── go.mod, go.sum
└── main.go, dev.go, prod.go
```

---

## Именование

| Что | Стиль | Пример |
|-----|-------|--------|
| Пакеты | `snake_case` + суффикс роли | `supplier_handler`, `supplier_service`, `order_dto`, `user_model` |
| Файлы | `snake_case` + суффикс роли | `supplier_handler.go`, `order_service.go`, `order_constants.go` |
| Экспортируемые типы / интерфейсы | `PascalCase` | `SupplierService`, `OrderCreate`, `ActivityStatus` |
| Приватные структуры | `camelCase` | `supplierService`, `supplierHandler` |
| Конструкторы | `New<Type>` | `NewSupplierService`, `NewSupplierHandler` |
| Методы | `PascalCase` | `FindByID`, `DeleteOrRestore`, `Page` |
| Переменные | `camelCase` | `supplierDto`, `supplierModel`, `authMiddleware` |
| Константы-статусы | `PascalCase` + описательный префикс | `StatusPending`, `PaymentTypeCash`, `PaymentStatusPaid` |
| Enum-значения | `SCREAMING_CASE` | `ACTIVE`, `INACTIVE`, `ASC`, `DESC` |
| Env-поля в struct | `SCREAMING_SNAKE` | `HTTP_Host`, `POSTGRES_DBName`, `JWT_Secret` |
| Алиасы импортов | `<package>_<role>` | `supplier_service "..."`, `auth_middleware "..."` |

---

## Паттерны

### Интерфейс + приватная реализация

```go
type SupplierService interface {
    Page(ctx context.Context, paginate *request.Paginate, filter pg.Filter) (*supplier_dto.SupplierPage, error)
    Find(ctx context.Context, filter pg.Filter) ([]supplier_dto.Supplier, error)
    Create(ctx context.Context, dto *supplier_dto.SupplierCreate) (int64, error)
    Update(ctx context.Context, id int64, dto *supplier_dto.SupplierUpdate) error
    DeleteOrRestore(ctx context.Context, id int64) error
}

type supplierService struct {
    db *gorm.DB
}

func NewSupplierService(db *gorm.DB) SupplierService {
    return &supplierService{db: db}
}
```

### Анонимный блок `{}` для группировки ошибок

Каждый вызов, возвращающий ошибку, оборачивается в `{}`:

```go
id, err := req.ParamToInt("id")
{
    if err != nil {
        return req.BadRequest(err)
    }
}

var supplierDto supplier_dto.SupplierUpdate
{
    if err := req.BindBody(&supplierDto); err != nil {
        return req.BadRequest(err)
    }
}
```

### cmd.go — точка входа модуля

```go
func Cmd(router *echo.Echo, db *gorm.DB, log logger.Logger, authMiddleware *auth_middleware.AuthMiddleware) {
    routerGroup := router.Group("/api/v1")
    {
        handler.NewXxxHandler(routerGroup, db, log, authMiddleware)
    }
}
```

### Регистрация маршрутов в handler

```go
func NewXxxHandler(router *echo.Echo, db *gorm.DB, log logger.Logger, authMiddleware *auth_middleware.AuthMiddleware) {
    handler := &xxxHandler{
        db:             db,
        log:            log,
        authMiddleware: authMiddleware,
        xxxService:     xxx_service.NewXxxService(db),
    }

    m := handler.authMiddleware.BuildMiddleware()
    group := router.Group("/api/v1/xxx", m)
    {
        group.POST("", handler.Create)
        group.PUT("/:id", handler.Update)
        group.GET("/:id", handler.FindByID)
        group.GET("/page", handler.Page)
        group.GET("/search", handler.Search)
        group.DELETE("/:id", handler.DeleteOrRestore)
    }
}
```

### Filter-функции для GORM

```go
filter := func(tx *gorm.DB) *gorm.DB {
    if params.Name != nil {
        tx = tx.Where("suppliers.name ILIKE ?", fmt.Sprintf("%%%s%%", *params.Name))
    }
    if params.OnlyDeleted != nil && *params.OnlyDeleted {
        tx = tx.Unscoped().Where("suppliers.deleted_at IS NOT NULL")
    } else if params.IncludeDeleted != nil && *params.IncludeDeleted {
        tx = tx.Unscoped()
    }
    return tx.Select("suppliers.id", "suppliers.name", "suppliers.created_at").
        Order("suppliers.id DESC")
}
```

### Мягкое удаление / восстановление (один метод на оба действия)

```go
func (s *xxxService) DeleteOrRestore(ctx context.Context, id int64) error {
    var entity Model

    if err := s.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&entity).Error; err != nil {
        return err
    }

    if entity.DeletedAt.Valid {
        return s.db.WithContext(ctx).Model(&Model{}).Where("id = ?", id).Update("deleted_at", nil).Error
    }
    return s.db.WithContext(ctx).Where("id = ?", id).Delete(&Model{}).Error
}
```

### Swagger godoc (обязателен для каждого публичного метода)

```go
// Create godoc
// @Summary      Create supplier
// @Description  Create new supplier
// @Tags         supplier
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        input body supplier_dto.SupplierCreate true "supplier information"
// @Success      201 {object} response.ID64 "Successful operation"
// @Failure      400 {object} response.HttpSuccess "Bad request"
// @Failure      500 {object} response.HttpSuccess "Internal server error"
// @Router       /supplier [POST]
func (h *supplierHandler) Create(c echo.Context) error {
```

### DTO — теги на каждом поле, аннотация `// @Name`

```go
type SupplierCreate struct {
    Name string `json:"name" validate:"required"`
} // @Name SupplierCreate

type SupplierUpdate struct {
    Name string `json:"name" validate:"required"`
} // @Name SupplierUpdate

type SupplierPage = response.PageData[Supplier] // @name SupplierPage

type SupplierParams struct {
    Name           *string `query:"name" json:"name"`
    IncludeDeleted *bool   `query:"include_deleted" json:"include_deleted"`
    OnlyDeleted    *bool   `query:"only_deleted" json:"only_deleted"`
} // @Name SupplierParams
```

### Model — минимальный GORM

```go
type Supplier struct {
    ID        int64           `json:"id" gorm:"primaryKey;autoIncrement"`
    Name      string          `json:"name" gorm:"type:varchar(255);not null"`
    CreatedAt *time.Time      `json:"created_at" gorm:"autoCreateTime"`
    DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`
} // @name Supplier
```

### Enum — тип + методы

```go
type ActivityStatus string

const (
    ACTIVE   ActivityStatus = "ACTIVE"
    INACTIVE ActivityStatus = "INACTIVE"
)

func (s ActivityStatus) IsValid() bool   { return s == ACTIVE || s == INACTIVE }
func (s ActivityStatus) Bool() bool      { return s == ACTIVE }
func (s ActivityStatus) String() string  { return string(s) }
// + MarshalJSON, UnmarshalJSON, Value, Scan — для корректной работы с JSON и БД
```

---

## Безопасность указателей и логика

### Nil-проверка перед разыменованием

Любое поле типа указатель, полученное из внешней библиотеки или из БД, проверяется перед использованием:

```go
// ❌
txInfo.From = intMsg.SrcAddr.String()

// ✅
if intMsg.SrcAddr != nil {
    txInfo.From = intMsg.SrcAddr.String()
}
```

### Все возвращаемые значения обязательны

Функция, возвращающая `(value, error)`, всегда обрабатывается с обоими значениями:

```go
// ❌ — теряем ошибку, возможна паника
payload := intMsg.Body.BeginParse()

// ✅
payload, err := intMsg.Body.BeginParse()
if err != nil { ... }

// ✅ — или внутри if
if payload, err := intMsg.Body.BeginParse(); err == nil {
    ...
}
```

### Type assertion — только на самом интерфейсе

```go
// ❌ — поле Description у any не существует
if desc, ok := tx.Description.Description.(tlb.TransactionDescriptionOrdinary); ok {

// ✅ — type assert напрямую на переменную
if desc, ok := tx.Description.(tlb.TransactionDescriptionOrdinary); ok {
    success = !desc.Aborted
}
```

### Статус результата берётся из реального состояния, не из предположения

```go
// ❌ — Success всегда true, реальные сбои не видны
txInfo := &TransactionInfo{Success: true}

// ✅ — берём из фактического результата операции
success := true
if desc, ok := tx.Description.(tlb.TransactionDescriptionOrdinary); ok {
    success = !desc.Aborted
}
txInfo := &TransactionInfo{Success: success}
```

### Агрегация в цикле — без преждевременного `break`

Если цикл обходит коллекцию для накопления значения — `break` запрещён, иначе остальные элементы игнорируются:

```go
// ❌ — суммируется только первый элемент
for _, msg := range list {
    txInfo.Amount = msg.Amount.TON()
    break
}

// ✅ — накапливаем всё
totalNano := new(big.Int)
for _, msg := range list {
    if msg.Amount.Nano() != nil {
        totalNano.Add(totalNano, msg.Amount.Nano())
    }
}
```

---

## Deprecated — запрещённые паттерны

### `ioutil` — удалён в Go 1.16+

Пакет `io/ioutil` полностью устарел. Все замены — в стандартной библиотеке:

```go
// ❌
ioutil.ReadAll(r)
ioutil.ReadFile(path)
ioutil.WriteFile(path, data, perm)
ioutil.TempFile(dir, pattern)
ioutil.TempDir(dir, pattern)
ioutil.Discard
ioutil.NopCloser(r)

// ✅
io.ReadAll(r)
os.ReadFile(path)
os.WriteFile(path, data, perm)
os.CreateTemp(dir, pattern)
os.MkdirTemp(dir, pattern)
io.Discard
io.NopCloser(r)
```

### HTTP-запросы — всегда с контекстом

```go
// ❌ — нет контекста, нельзя отменить, нет таймаута
resp, err := http.Get(url)
resp, err := http.Post(url, contentType, body)

// ✅
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
if err != nil { return err }
resp, err := http.DefaultClient.Do(req)
```

### Ошибки — оборачивать через `%w`

```go
// ❌ — потеря контекста, нельзя проверить через errors.Is/As
return errors.New("failed: " + err.Error())
return fmt.Errorf("failed: " + err.Error())

// ✅
return fmt.Errorf("failed to create wallet: %w", err)
```

### Проверка ошибок файловой системы

```go
// ❌ — устаревший os.IsNotExist
if os.IsNotExist(err) { ... }

// ✅
if errors.Is(err, os.ErrNotExist) { ... }
```

### Генератор случайных чисел — `rand.Seed` удалён в Go 1.20

```go
// ❌
rand.Seed(time.Now().UnixNano())
n := rand.Intn(100)

// ✅ — глобальный rand автоматически инициализирован с Go 1.20
n := rand.IntN(100)  // Go 1.22+
// или
src := rand.NewSource(time.Now().UnixNano())
r := rand.New(src)
n := r.Intn(100)
```

### GORM — устаревшие паттерны

```go
// ❌ — First с условием строкой (нарушает явность)
db.First(&model, "name = ?", name)

// ✅ — явный Where
db.Where("name = ?", name).First(&model)

// ❌ — Find без Select (загружает все колонки)
db.Find(&models)

// ✅ — явный Select нужных полей
db.Select("id", "name", "created_at").Find(&models)

// ❌ — db.Table("name") вместо модели
db.Table("suppliers").Where(...)

// ✅ — через типизированную модель
db.Model(&Supplier{}).Where(...)
```

### Context — обязательный первый аргумент

```go
// ❌ — context.Background() в методах сервиса
func (s *service) Create() error {
    s.db.WithContext(context.Background()).Create(...)
}

// ✅ — ctx приходит снаружи и пробрасывается
func (s *service) Create(ctx context.Context) error {
    s.db.WithContext(ctx).Create(...)
}
```

---

## Запрещено

- ❌ `panic` — только возвращать ошибки
- ❌ Глобальные переменные — всё через DI конструктора
- ❌ Бизнес-логика в handler — только вызов сервиса
- ❌ Смешивать DTO и model
- ❌ Пропускать Swagger аннотации на публичных методах
- ❌ `Select("*")` — явно перечислять нужные поля
- ❌ `interface{}` или `any` в публичных сигнатурах
- ❌ Комментарии в коде (исключения: `// TODO:` и Swagger godoc)
- ❌ Разыменовывать указатель без проверки на `nil`
- ❌ Игнорировать возвращаемую `error` (присваивать только первое значение `_`)
- ❌ `break` в цикле агрегации — если цель накопить, обойди всё
- ❌ Хардкодить статус результата (`Success: true`) — брать из реального состояния
- ❌ `io/ioutil` — использовать `io` и `os`
- ❌ `http.Get/Post` без контекста — только `http.NewRequestWithContext`
- ❌ `errors.New(err.Error())` — оборачивать через `fmt.Errorf("%w", err)`
- ❌ `os.IsNotExist` — использовать `errors.Is(err, os.ErrNotExist)`
- ❌ `rand.Seed` — не нужен с Go 1.20+
- ❌ `db.Table("name")` — только через типизированную модель `db.Model(&T{})`
- ❌ `context.Background()` внутри метода сервиса — ctx всегда приходит параметром

---

## Комментарии

Комментарии **не пишутся**. Единственные исключения:
- `// TODO:` для временных пометок
- Swagger godoc аннотации (обязательная документация API)

Если комментарий всё-таки нужен — писать **по-русски**.

---

## Тесты

- Тесты только на бэкенде, фронтенд не тестируется
- Папка `idoctor_back/tests/` — на одном уровне с `src/`

```
idoctor_back/
├── src/
└── tests/
    └── <name>_service_test.go
```

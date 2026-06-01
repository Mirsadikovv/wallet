# CLAUDE_STYLE — Telegram Bot (Go)

> Правила для кода в `idoctor_bot/`. Строго следуй при написании нового кода.

---

## Стек

- **Язык**: Go
- **Библиотека**: go-telegram-bot-api/telegram-bot-api/v5
- **Режим**: Long polling (`GetUpdatesChan`) — Webhook не используется
- **БД**: PostgreSQL через GORM (shared с бэкендом)
- **Config**: `godotenv` + `os.Getenv`
- **Shared**: `github.com/Mirsadikovv/shared/pg` для подключения к БД

---

## Структура

```
idoctor_bot/
├── app/
│   ├── bot/
│   │   └── bot.go              ← запуск бота, update loop
│   ├── config/
│   │   └── config.go           ← Config struct + Load()
│   ├── handler/
│   │   └── privates/
│   │       ├── handlers.go     ← Register() + HandleUpdate() (главный роутер)
│   │       ├── texts.go        ← все мультиязычные тексты (map[string]string)
│   │       ├── start.go        ← обработчик /start
│   │       ├── menu.go         ← главное меню
│   │       ├── back.go         ← кнопка "назад"
│   │       ├── appeal.go       ← обработчик заявки
│   │       ├── change_language.go ← выбор и смена языка
│   │       ├── send_phone_number.go ← получение номера телефона
│   │       ├── help.go         ← /help
│   │       └── echo.go         ← дефолтный обработчик (fallback)
│   ├── keyboards/
│   │   └── defaults/
│   │       ├── main_menu.go    ← главное меню
│   │       ├── languages.go    ← выбор языка
│   │       ├── send_contact.go ← кнопка поделиться контактом
│   │       └── appeal.go       ← клавиатура заявки
│   ├── utils/
│   │   ├── build_keyboard.go   ← MakeReplyMarkup
│   │   ├── commands.go         ← GetBotCommands
│   │   ├── current_lang.go     ← LanguageCache
│   │   └── notify_admin.go     ← NotifyAdmins
│   └── main.go                 ← Exec(env *Env)
├── main.go                     ← точка входа: main() + fillEnv()
├── dev.go                      ← заполнение env для dev
└── prod.go                     ← заполнение env для prod
```

**Правило**: каждый обработчик — отдельный файл. Имя файла = имя функции в `snake_case`.

---

## Именование

| Что | Стиль | Пример |
|-----|-------|--------|
| Пакеты | короткое имя без суффикса | `bot`, `config`, `handlers`, `keyboard`, `utils` |
| Файлы | `snake_case` | `change_language.go`, `send_phone_number.go`, `build_keyboard.go` |
| Экспортируемые функции | `PascalCase` | `Start`, `HandleUpdate`, `Register`, `ChooseLanguage`, `MakeReplyMarkup` |
| Приватные функции | `camelCase` | `clauseOnConflict` |
| Тексты (карты переводов) | `PascalCase` | `StartWriteYourPhone`, `MenuText`, `HelpText`, `Back` |
| Клавиатуры (JSON-константы) | `PascalCase` | `LanguageKeyboard`, `SendContact`, `MainMenuKeyboard` |
| Карты клавиатур по языкам | `PascalCase` + `Map` | `SendContactMap`, `AppealKeyboardSend` |
| Config поля | `PascalCase` | `BotToken`, `AdminIds`, `Debug` |
| Env-поля в struct | `SCREAMING_SNAKE` | `POSTGRES_HOST`, `POSTGRES_DB` |

---

## Паттерны

### Update loop — bot.go

```go
func Start(cfg *config.Config, db *gorm.DB) error {
    bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
    {
        if err != nil {
            log.Printf("failed to create bot: %v", err)
            return err
        }
    }

    bot.Debug = false

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 100

    handlers.Register(bot, db)
    utils.NotifyAdmins(bot, cfg)

    updates := bot.GetUpdatesChan(u)
    langCache := utils.NewLanguageCache()

    for update := range updates {
        if update.Message != nil {
            curLang := handlers.HandleUpdate(bot, update, cfg, db, langCache)
            langCache.Set(update.Message.From.ID, curLang.Get(update.Message.From.ID))
        }
    }

    return nil
}
```

### Главный роутер — HandleUpdate

Порядок проверок строгий: сначала `Contact`, потом команды, потом текст:

```go
func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, db *gorm.DB, langCache *utils.LanguageCache) *utils.LanguageCache {

    if update.Message.From.LanguageCode != "" && langCache.Get(update.Message.From.ID) == "" {
        langCache.Set(update.Message.From.ID, update.Message.From.LanguageCode)
    }

    if update.Message.Contact != nil && update.Message.Contact.PhoneNumber != "" {
        return SendPhoneNumber(bot, update, db, langCache)

    } else if update.Message != nil && update.Message.IsCommand() {
        switch update.Message.Command() {
        case "start":
            return Start(bot, update, cfg, db, langCache)
        case "help":
            Echo(bot, update, cfg, langCache)
        default:
            Echo(bot, update, cfg, langCache)
        }
    } else if update.Message != nil {
        switch update.Message.Text {
        case "Назад ⬅️", "Ortga ⬅️", "Back ⬅️":
            return Back(bot, update, cfg, db, langCache)
        case "Change language🌐", "Tilni o'zgartirish🌐", "Поменять язык🌐":
            return ChooseLanguage(bot, update, langCache)
        case "🇬🇧 English", "🇷🇺 Русский", "🇺🇿 O'zbek":
            return ChangeLanguage(bot, update, db, langCache)
        default:
            Echo(bot, update, cfg, langCache)
        }
    }
    return langCache
}
```

### Сигнатура обработчика

Все обработчики принимают `langCache` и возвращают его (даже если не меняют):

```go
func HandlerName(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, db *gorm.DB, lang *utils.LanguageCache) *utils.LanguageCache {
    // ...
    return lang
}
```

Простые обработчики без изменения языка могут не возвращать:

```go
func Appeal(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    // ...
}
```

### Отправка сообщения

```go
msg := tgbotapi.NewMessage(update.Message.Chat.ID, TextMap[lang.Get(update.Message.From.ID)])
msg.ParseMode = "HTML"
msg.ReplyMarkup = utils.MakeReplyMarkup(keyboard.SomeKeyboard)

if _, err := bot.Send(msg); err != nil {
    log.Println("Error sending message:", err)
}
```

### Тексты — все в texts.go, map[string]string по языкам

```go
var MenuText = map[string]string{
    "uz": "Xush kelibsiz, <a href=\"tg://user?id=%d\">%s</a>!",
    "ru": "Добро пожаловать, <a href=\"tg://user?id=%d\">%s</a>!",
    "en": "Welcome, <a href=\"tg://user?id=%d\">%s</a>!",
}

var BackChooseAction = map[string]string{
    "ru": "Выберите действие:",
    "uz": "Tanlang:",
    "en": "Choose action:",
}
```

Ключи языков: `"uz"`, `"ru"`, `"en"` — строго три языка в таком порядке. Все новые тексты — в `texts.go`.

### Клавиатуры — raw JSON строки

Клавиатуры описываются как JSON-строки в `const` или `var`, с `fmt.Sprintf` для динамических данных:

```go
const MainMenuKeyboard = `
    {
        "keyboard": [
        [
            {
                "text": "%s",
                "web_app": { "url": "https://example.uz" }
            }
        ],
        [
            {
                "text": "%s",
                "callback_data": "language"
            }
        ]
        ],
        "resize_keyboard": true,
        "one_time_keyboard": true
    }`

func GetMainMenuKeyboard(lang string, userID int64) string {
    return fmt.Sprintf(MainMenuKeyboard, ButtonText1[lang], ButtonText2[lang])
}
```

Карты кнопок по языкам — рядом с клавиатурой в том же файле:

```go
var (
    ButtonLabel = map[string]string{
        "ru": "Текст кнопки",
        "uz": "Tugma matni",
        "en": "Button text",
    }
)
```

### MakeReplyMarkup — преобразование JSON-строки в разметку

```go
func MakeReplyMarkup(jsonStr string) json.RawMessage {
    var raw json.RawMessage
    err := json.Unmarshal([]byte(jsonStr), &raw)
    if err != nil {
        log.Println("Error in MakeReplyMarkup:", err)
        return nil
    }
    return raw
}
```

Использование: `msg.ReplyMarkup = utils.MakeReplyMarkup(keyboard.SomeKeyboard)`

### LanguageCache — in-memory кеш языка пользователя

```go
type LanguageCache struct {
    m map[int64]string
}

func NewLanguageCache() *LanguageCache {
    return &LanguageCache{m: make(map[int64]string)}
}

func (lc *LanguageCache) Set(userID int64, language string) { lc.m[userID] = language }
func (lc *LanguageCache) Get(userID int64) string           { return lc.m[userID] }
```

Язык берётся: `lang.Get(update.Message.From.ID)`.

### Уведомление администраторов

```go
func NotifyAdmins(bot *tgbotapi.BotAPI, cfg *config.Config) {
    for _, idStr := range cfg.AdminIds {
        id, err := strconv.ParseInt(idStr, 10, 64)
        if err != nil {
            log.Printf("Invalid admin ID: %v", err)
            continue
        }
        // send message...
    }
}
```

Паттерн рассылки по всем админам с `continue` при ошибке парсинга ID.

### Config — через godotenv + os.Getenv

```go
type Config struct {
    BotToken string
    AdminIds []string
    Debug    bool
}

func Load() (*Config, error) {
    _ = godotenv.Load(".env")
    adminEnv := os.Getenv("ADMINS")
    admins := strings.Split(adminEnv, ",")
    return &Config{
        BotToken: os.Getenv("BOT_TOKEN"),
        AdminIds: admins,
    }, nil
}
```

`ADMINS` — через запятую: `"123456,789012"`.

### Работа с БД в хендлерах — напрямую через db

```go
if err := db.Table("users").
    Where("telegram_id = ?", update.Message.From.ID).
    Update("phone_number", update.Message.Contact.PhoneNumber).Error; err != nil {
    log.Println("Error updating bot user state:", err)
}
```

Нет сервисного слоя — `db` передаётся напрямую в хендлер и используется там же. Это допустимо для бота.

### Upsert через OnConflict

```go
func clauseOnConflict() clause.OnConflict {
    return clause.OnConflict{
        Columns: []clause.Column{{Name: "telegram_id"}},
        DoUpdates: clause.Assignments(map[string]any{
            "telegram_id":       gorm.Expr("excluded.telegram_id"),
            "telegram_username": gorm.Expr("excluded.telegram_username"),
        }),
    }
}

db.Table("users").Clauses(clauseOnConflict()).Create(&newUser)
```

### Ошибки — только log, без return

В хендлерах ошибки не возвращаются, только логируются:

```go
if _, err := bot.Send(msg); err != nil {
    log.Println("Error sending message:", err)
}
```

---

## Запрещено

- ❌ Webhook — только long polling
- ❌ Inline keyboards / callback queries — только reply keyboards
- ❌ Goroutine на каждый update
- ❌ FSM/state machine — только switch по тексту сообщения
- ❌ Сервисный слой — `db` используется напрямую в хендлерах
- ❌ Структурированное логирование — только `log.Println` / `log.Printf`
- ❌ Комментарии в коде (исключение: `// TODO:`)
- ❌ Новые тексты вне `texts.go` — все переводы только там

---

## Добавление нового обработчика

1. Создать файл `app/handler/privates/<name>.go`
2. Добавить тексты в `app/handler/privates/texts.go`
3. Если нужна клавиатура — создать `app/keyboards/defaults/<name>.go`
4. Зарегистрировать в `HandleUpdate` в `handlers.go`

```go
// 1. app/handler/privates/my_feature.go
package handlers

import (
    keyboard "github.com/Mirsadikovv/idoctor_bot/app/keyboards/defaults"
    "github.com/Mirsadikovv/idoctor_bot/app/utils"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MyFeature(bot *tgbotapi.BotAPI, update tgbotapi.Update, lang *utils.LanguageCache) *utils.LanguageCache {
    userLang := lang.Get(update.Message.From.ID)

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, MyFeatureText[userLang])
    msg.ParseMode = "HTML"
    msg.ReplyMarkup = utils.MakeReplyMarkup(keyboard.MyFeatureKeyboard)

    if _, err := bot.Send(msg); err != nil {
        log.Println("Error sending message:", err)
    }

    return lang
}

// 2. texts.go — добавить
var MyFeatureText = map[string]string{
    "uz": "...",
    "ru": "...",
    "en": "...",
}

// 3. keyboards/defaults/my_feature.go
package keyboard

const MyFeatureKeyboard = `{ "keyboard": [...], "resize_keyboard": true }`

// 4. handlers.go — добавить в switch
case "Кнопка RU", "Tugma UZ", "Button EN":
    return MyFeature(bot, update, lang)
```

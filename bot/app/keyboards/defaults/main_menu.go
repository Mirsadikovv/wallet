package keyboard

import "fmt"

var OpenWalletButton = map[string]string{
	"uz": "💼 Hamyonni ochish",
	"ru": "💼 Открыть кошелёк",
	"en": "💼 Open Wallet",
}

var ChangeLangButton = map[string]string{
	"uz": "Tilni o'zgartirish🌐",
	"ru": "Поменять язык🌐",
	"en": "Change language🌐",
}

const mainMenuKeyboardTpl = `{
    "keyboard": [
        [{"text": "%s", "web_app": {"url": "%s"}}],
        [{"text": "%s"}]
    ],
    "resize_keyboard": true
}`

func GetMainMenuKeyboard(lang, webAppURL string) string {
	return fmt.Sprintf(mainMenuKeyboardTpl, OpenWalletButton[lang], webAppURL, ChangeLangButton[lang])
}

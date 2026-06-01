package keyboard

import "fmt"

var BackButton = map[string]string{
	"uz": "Ortga ⬅️",
	"ru": "Назад ⬅️",
	"en": "Back ⬅️",
}

const backKeyboardTpl = `{
    "keyboard": [
        [{"text": "%s"}]
    ],
    "resize_keyboard": true
}`

func GetBackKeyboard(lang string) string {
	return fmt.Sprintf(backKeyboardTpl, BackButton[lang])
}

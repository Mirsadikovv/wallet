package handlers

var WelcomeText = map[string]string{
	"uz": "Xush kelibsiz, <a href=\"tg://user?id=%d\">%s</a>!\n\nIltimos, elektron pochtangizni kiriting:",
	"ru": "Добро пожаловать, <a href=\"tg://user?id=%d\">%s</a>!\n\nПожалуйста, введите ваш email:",
	"en": "Welcome, <a href=\"tg://user?id=%d\">%s</a>!\n\nPlease enter your email:",
}

var MenuText = map[string]string{
	"uz": "Xush kelibsiz, <a href=\"tg://user?id=%d\">%s</a>! 👋",
	"ru": "Добро пожаловать, <a href=\"tg://user?id=%d\">%s</a>! 👋",
	"en": "Welcome, <a href=\"tg://user?id=%d\">%s</a>! 👋",
}

var EnterEmailText = map[string]string{
	"uz": "Elektron pochtangizni kiriting:",
	"ru": "Введите ваш email:",
	"en": "Enter your email:",
}

var OTPSentText = map[string]string{
	"uz": "📧 <b>%s</b> manziliga tasdiqlash kodi yuborildi.\n\nKodni kiriting:",
	"ru": "📧 Код подтверждения отправлен на <b>%s</b>.\n\nВведите код:",
	"en": "📧 Verification code sent to <b>%s</b>.\n\nEnter the code:",
}

var OTPSuccessText = map[string]string{
	"uz": "✅ Muvaffaqiyatli ro'yxatdan o'tdingiz! Hamyoningiz yaratildi.",
	"ru": "✅ Регистрация прошла успешно! Кошелёк создан.",
	"en": "✅ Registration successful! Your wallet has been created.",
}

var OTPInvalidText = map[string]string{
	"uz": "❌ Noto'g'ri yoki muddati o'tgan kod. Qayta urinib ko'ring.",
	"ru": "❌ Неверный или просроченный код. Попробуйте снова.",
	"en": "❌ Invalid or expired code. Please try again.",
}

var InvalidEmailText = map[string]string{
	"uz": "❌ Noto'g'ri email format. Iltimos, to'g'ri email kiriting:",
	"ru": "❌ Неверный формат email. Пожалуйста, введите корректный email:",
	"en": "❌ Invalid email format. Please enter a valid email:",
}

var EmailTakenText = map[string]string{
	"uz": "❌ Bu email allaqachon ro'yxatdan o'tgan. Boshqa email kiriting:",
	"ru": "❌ Этот email уже зарегистрирован. Введите другой email:",
	"en": "❌ This email is already registered. Please enter a different email:",
}

var ChooseLanguageText = map[string]string{
	"uz": "Tilni tanlang:",
	"ru": "Выберите язык:",
	"en": "Choose language:",
}

var HelpText = map[string]string{
	"uz": "📱 <b>TON Hamyon</b>\n\nHamyonni boshqarish uchun <b>Hamyonni ochish</b> tugmasini bosing.",
	"ru": "📱 <b>TON Кошелёк</b>\n\nДля управления кошельком нажмите кнопку <b>Открыть кошелёк</b>.",
	"en": "📱 <b>TON Wallet</b>\n\nPress <b>Open Wallet</b> button to manage your wallet.",
}

var EchoText = map[string]string{
	"uz": "Iltimos, menyu tugmalaridan foydalaning.",
	"ru": "Пожалуйста, используйте кнопки меню.",
	"en": "Please use the menu buttons.",
}

var RegistrationErrorText = map[string]string{
	"uz": "❌ Ro'yxatdan o'tishda xatolik yuz berdi. Qayta urinib ko'ring /start",
	"ru": "❌ Ошибка при регистрации. Попробуйте снова /start",
	"en": "❌ Registration failed. Please try again /start",
}

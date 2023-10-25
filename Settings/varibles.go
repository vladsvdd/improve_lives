package Settings

import (
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	//tapioca "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var (
	GlobalPort string
	// GlobalBot Бот Telegram
	GlobalBot *tapioca.BotAPI
	// GlobalToken Токен бота Telegram
	GlobalToken string
	// GlobalDatabase Соединение с базой данных
	GlobalDatabase *gorm.DB

	// GlobalTest запускаем телеграмм бота в тестовом окружении
	GlobalTest bool

	GlobalDbUsername string
	GlobalDbPassword string
	GlobalDbHost     string
	GlobalDbPort     string
	GlobalDbDatabase string
)

package Core

import (
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	"improve_lives/db/objects"
	"log"
	"strings"
)

// StoreUserFromUpdate сохраняет пользователя из update
func StoreUserFromUpdate(update *tapioca.Update) (user *objects.User) {
	if CallbackQueryFromIsMissing(update) {
		log.Printf("0|||||Error in StoreUserFromUpdate() -> CallbackQueryFromIsMissing() is false")
		return
	}

	from := update.CallbackQuery.From

	// Формируем структуру User
	user = &objects.User{
		ChatID:       int64(from.ID),
		IsBot:        from.IsBot,
		UserName:     strings.TrimSpace(from.UserName),
		FirstName:    strings.TrimSpace(from.FirstName),
		LastName:     strings.TrimSpace(from.LastName),
		LanguageCode: strings.TrimSpace(from.LanguageCode),
	}

	return user
}

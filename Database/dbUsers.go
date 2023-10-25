package Database

import (
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	//tapioca "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
	"strings"
)

// StoreUserToDB Функция для записи пользователя в базу данных
func StoreUserToDB(user *objects.User) (err error) {
	// Проверяем, существует ли запись с указанным chat_id
	err = Settings.GlobalDatabase.Table(objects.User{}.TableName()).Where("chat_id = ?", user.ChatID).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		if err = Settings.GlobalDatabase.Table(objects.User{}.TableName()).Create(&user).Error; err != nil {
			return err
		}
		InsertStartActivitiesAndRewards(user)
	}

	return err
}

// GetUserFromDB Функция для получения пользователя из базы данных
func GetUserFromDB(update *tapioca.Update) (user *objects.User, err error) {
	// Подготавливаем SQL-запрос для выборки данных пользователя
	err = Settings.GlobalDatabase.Table(objects.User{}.TableName()).Where("chat_id = ?", update.Message.From.ID).First(&user).Error

	// Если пользователя нет в базе данных, создайте нового пользователя
	if err != nil {
		user = &objects.User{
			ChatID:       int64(update.Message.From.ID),
			IsBot:        update.Message.From.IsBot,
			UserName:     strings.TrimSpace(update.Message.From.UserName),
			FirstName:    strings.TrimSpace(update.Message.From.FirstName),
			LastName:     strings.TrimSpace(update.Message.From.LastName),
			LanguageCode: strings.TrimSpace(update.Message.From.LanguageCode),
		}
		// Сохраните нового пользователя в базу данных
		err = StoreUserToDB(user)
		if err != nil {
			log.Printf("%d|||||%s() Данные пользователя не записаны в БД. %s", user.ChatID, "GetUserFromDB_1", err)
		}

		InsertStartActivitiesAndRewards(user)
	}

	return user, err
}

func UpdateUserActionIdFromDB(user *objects.User, update *tapioca.Update) {
	choiceCode := update.CallbackQuery.Data
	userAction := objects.UserAction{}

	err := Settings.GlobalDatabase.Table(objects.UserAction{}.TableName()).Where("action = ?", choiceCode).First(&userAction).Error
	if err != nil {
		userAction.Id = 0
	}

	err = Settings.GlobalDatabase.Model(&user).Update("user_action_id", userAction.Id).Error
	if err != nil {
		log.Printf("%s() %v", "UpdateUserActionIdFromDB", err)
	}
}

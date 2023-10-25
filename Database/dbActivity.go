package Database

import (
	"github.com/google/uuid"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
)

func GetActivitiesByUser(user *objects.User) (activities []objects.Activity) {
	// Выполняем запрос с использованием GORM
	Settings.GlobalDatabase.Table(objects.Activity{}.TableName()). // .Select(...)
									Joins("LEFT JOIN user ON activity.user_id = user.id").
									Where("user.chat_id = ? AND activity.is_deleted = 0", user.ChatID).
									Find(&activities)

	if Settings.GlobalDatabase.Error != nil {
		log.Printf("Ошибка при выполнении запроса: %v", Settings.GlobalDatabase.Error)
	}

	// В переменной activities теперь содержится массив с данными из БД
	return activities
}

func AddNewActivity(user *objects.User, coins int, text string) error {
	// Генерируем новый UUID
	id := uuid.New()
	// Преобразуем UUID в строку
	uniqueString := id.String()

	activity := objects.Activity{
		UserId: user.Id,
		Coins:  int64(coins),
		Name:   text,
		Code:   uniqueString,
	}

	if err := Settings.GlobalDatabase.Table(objects.Activity{}.TableName()).Create(&activity).Error; err != nil {
		return err
	}

	return nil
}

package Database

import (
	"github.com/google/uuid"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
)

func GetRewardsByUser(user *objects.User) (rewards []objects.Reward) {
	// Выполняем запрос с использованием GORM
	Settings.GlobalDatabase.Table(objects.Reward{}.TableName()). // .Select(...)
									Joins("LEFT JOIN user ON reward.user_id = user.id").
									Where("user.chat_id = ? AND is_deleted = 0", user.ChatID).
									Find(&rewards)

	if Settings.GlobalDatabase.Error != nil {
		log.Printf("Ошибка при выполнении запроса: %v", Settings.GlobalDatabase.Error)
	}
	// В переменной activities теперь содержится массив с данными из БД
	return rewards
}

func AddNewReward(user *objects.User, coins int, text string) error {
	// Генерируем новый UUID
	id := uuid.New()
	// Преобразуем UUID в строку
	uniqueString := id.String()

	reward := objects.Reward{
		UserId: user.Id,
		Coins:  int64(coins),
		Name:   text,
		Code:   uniqueString,
	}

	if err := Settings.GlobalDatabase.Table(objects.Reward{}.TableName()).Create(&reward).Error; err != nil {
		return err
	}

	return nil
}

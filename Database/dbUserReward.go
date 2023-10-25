package Database

import (
	"improve_lives/Settings"
	"improve_lives/db/objects"
)

// StoreRewardToDB Функция для записи активности в базу данных
func StoreRewardToDB(reward *objects.Reward, user *objects.User) error {
	// Создаем связанные записи в таблице user_activities
	userReward := &objects.UserReward{
		UserId:       user.Id,
		UserRewardId: reward.Id,
	}

	if err := Settings.GlobalDatabase.Table(objects.UserReward{}.TableName()).Create(userReward).Error; err != nil {
		return err
	}

	return nil
}

package Database

import (
	"improve_lives/Settings"
	"improve_lives/db/objects"
)

// StoreActivityToDB Функция для записи активности в базу данных
func StoreActivityToDB(activity *objects.Activity, user *objects.User) error {
	// Создаем связанные записи в таблице user_activities
	userActivity := &objects.UserActivity{
		UserId:         user.Id,
		UserActivityId: activity.Id,
	}

	if err := Settings.GlobalDatabase.Table(objects.UserActivity{}.TableName()).Create(userActivity).Error; err != nil {
		return err
	}

	return nil
}

// GetSumCoinsActivityFromDB Функция для получения активности из базы данных
func GetSumCoinsActivityFromDB(user *objects.User) (int64, error) {
	var sumActivities int64
	var sumRewards int64
	query := `
		SELECT 
		  COALESCE(
		    (
		      SELECT 
		        SUM(a.coins) 
		      FROM 
		        user_activity AS ua 
		        LEFT JOIN activity AS a ON ua.user_activity_id = a.id 
		      WHERE 
		        a.user_id = u.id
		    ), 
		    0
		  ) AS sum, 
		  COALESCE(
		    (
		      SELECT 
		        SUM(coins) 
		      FROM 
		        user_reward AS ur 
		        LEFT JOIN reward AS r ON ur.user_reward_id = r.id 
		      WHERE 
		        r.user_id = u.id
		    ), 
		    0
		  ) AS sum2 
		FROM 
		  user AS u 
		WHERE 
		  u.chat_id = ?
	`

	//Ручной запрос в БД
	err := Settings.GlobalDatabase.Raw(query, user.ChatID).Row().Scan(&sumActivities, &sumRewards)
	if err != nil {
		return 0, err
	}

	return sumActivities - sumRewards, nil
}

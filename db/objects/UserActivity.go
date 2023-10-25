package objects

import "time"

// UserActivity представляет активность
type UserActivity struct {
	Id             int64 `gorm:"primaryKey"`
	UserId         int64 `gorm:"foreignKey:UserId;references:Id"`
	UserActivityId int64 `gorm:"foreignKey:UserActivityId;references:Id"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TableName Здесь мы указываем явное имя таблицы "user" для структуры User.
func (UserActivity) TableName() string {
	return "user_activity"
}

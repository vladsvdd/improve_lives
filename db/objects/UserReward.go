package objects

import "time"

// UserReward представляет активность
type UserReward struct {
	Id           int64 `gorm:"primaryKey"`
	UserId       int64 `gorm:"foreignKey:UserId;references:Id"`
	UserRewardId int64 `gorm:"foreignKey:UserRewardId;references:Id"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName Здесь мы указываем явное имя таблицы "user" для структуры User.
func (UserReward) TableName() string {
	return "user_reward"
}

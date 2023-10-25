package objects

import "time"

// Reward представляет награду
type Reward struct {
	Id        int64 `gorm:"primaryKey"`
	Code      string
	Name      string
	IsDeleted bool
	Coins     int64
	UserId    int64 `gorm:"foreignKey:UserId;references:Id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName Здесь мы указываем явное имя таблицы "user" для структуры User.
func (Reward) TableName() string {
	return "reward"
}

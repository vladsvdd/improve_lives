package objects

import "time"

// Activity представляет активность
type Activity struct {
	Id        int64 `gorm:"primaryKey"`
	Code      string
	Name      string
	IsDeleted bool
	Coins     int64
	UserId    int64 `gorm:"foreignKey:UserId;references:Id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName Здесь мы указываем явное имя таблицы "activity" для структуры Activity.
func (Activity) TableName() string {
	return "activity"
}

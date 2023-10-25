package objects

import "time"

type User struct {
	Id           int64 `gorm:"primaryKey"`
	ChatID       int64
	UserActionId int64 `gorm:"foreignKey:UserActionId;references:Id"`
	IsBot        bool
	UserName     string
	FirstName    string
	LastName     string
	LanguageCode string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName Здесь мы указываем явное имя таблицы "user" для структуры User.
func (User) TableName() string {
	return "user"
}

package objects

import "time"

type UserAction struct {
	Id        int64 `gorm:"primaryKey"`
	ParentId  int64
	Action    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName Здесь мы указываем явное имя таблицы "user" для структуры UserAction.
func (UserAction) TableName() string {
	return "user_action"
}

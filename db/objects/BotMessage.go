package objects

import "time"

// BotMessage представляет сообщения от бота
type BotMessage struct {
	Id        int64 `gorm:"primaryKey"`
	UserId    int64 `gorm:"foreignKey:UserId;references:Id"`
	MessageId int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName Здесь мы указываем явное имя таблицы "user" для структуры BotMessage.
func (BotMessage) TableName() string {
	return "bot_message"
}

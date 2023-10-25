package Database

import (
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	//tapioca "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
)

func CreateBotMessage(user *objects.User, message *tapioca.Message) (botMessage objects.BotMessage) {
	botMessage = objects.BotMessage{
		UserId:    user.Id,
		MessageId: message.MessageID,
	}

	if err := Settings.GlobalDatabase.Table(objects.BotMessage{}.TableName()).Create(&botMessage).Error; err != nil {
		log.Printf("dbBotMessage.go_api %s() %v", "create", err)
	}

	return botMessage
}

func GetRecordBotMessage(user *objects.User) (botMessage objects.BotMessage) {
	if err := Settings.GlobalDatabase.Table(objects.BotMessage{}.TableName()).Where("user_id = ?", user.Id).Last(&botMessage).Error; err != nil {
		log.Printf("dbBotMessage.go_api %s() %v", "getRecord", err)
	}
	return botMessage
}

func DeleteBotMessage(user *objects.User) {
	err := Settings.GlobalDatabase.
		Where("user_id = ?", user.Id).
		Delete(&objects.BotMessage{}).
		Error

	if err != nil {
		log.Printf("Ошибка при удалении сообщений из БД: %v", err)
	} else {
		log.Printf("Все сообщения пользователя успешно удалены из БД! UserId=%d", user.Id)
	}
}

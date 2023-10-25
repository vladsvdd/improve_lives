// Package Core
// bot.go_api - Файл, где вы можете определить функции для работы с ботом, отправки сообщений, кнопок и так далее.
// Функции, связанные с обработкой сообщений и логикой бота
package Core

import (
	"fmt"
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	"improve_lives/Core/scope"

	//tapioca "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"improve_lives/Database"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
	"time"
)

// IsStartMessage возвращает true, если сообщение - это команда /start
func IsStartMessage(update *tapioca.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

// IsCallbackQuery возвращает true, если сообщение - это ответ на inline-кнопку
func IsCallbackQuery(update *tapioca.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data != ""
}

// Delay приостанавливает выполнение программы на указанное количество секунд
func Delay(seconds uint8) {
	time.Sleep(time.Second * time.Duration(seconds))
}

// SendBotIsTyping отправить действие "печати"
func SendBotIsTyping(user *objects.User) {
	// Отправить действие "печати"
	chatAction := tapioca.NewChatAction(user.ChatID, tapioca.ChatTyping)
	_, err := Settings.GlobalBot.Send(chatAction)
	if err != nil {
		log.Printf("%d|||||%s Ошибка отправки действия печати:", user.ChatID, err)
	}
}

// SendStringMessage отправляет текстовое сообщение пользователю
func SendStringMessage(user *objects.User, msg string) {
	_, err := Settings.GlobalBot.Send(tapioca.NewMessage(user.ChatID, msg))
	if err != nil {
		log.Printf("%d||||| ERROR in sendStringMessage() %s", user.ChatID, err)
	}
}

// SendMessageWithDelay отправляет текстовое сообщение с задержкой в секундах
func SendMessageWithDelay(user *objects.User, delayInSec uint8, message string) {
	SendStringMessage(user, message)
	SendBotIsTyping(user)
	Delay(delayInSec)
}

// GetKeyboardRow возвращает строку с inline-кнопкой
func GetKeyboardRow(buttonText, buttonCode string) []tapioca.InlineKeyboardButton {
	return tapioca.NewInlineKeyboardRow(tapioca.NewInlineKeyboardButtonData(buttonText, buttonCode))
}

// CallbackQueryFromIsMissing возвращает true, если callback-запрос не содержит информации о пользователе
func CallbackQueryFromIsMissing(update *tapioca.Update) bool {
	return update.CallbackQuery == nil || update.CallbackQuery.From == nil
}

// SendUpdatedMessage редактированние предыдущего сообщения и отправка нового
func SendUpdatedMessage(
	user *objects.User,
	update *tapioca.Update,
	buttons tapioca.InlineKeyboardMarkup,
	message string,
	callingFunctionName string,
) {
	// Создайте объект tapioca.EditMessageTextConfig для обновления сообщения
	msg := tapioca.EditMessageTextConfig{
		BaseEdit: tapioca.BaseEdit{
			ChatID:      user.ChatID,
			MessageID:   update.CallbackQuery.Message.MessageID,
			ReplyMarkup: &buttons,
		},
		Text: message,
	}
	//msg.ReplyMarkup = tapioca.NewInlineKeyboardMarkup(activitiesButtonsRows...)
	_, err := Settings.GlobalBot.Send(msg)
	if err != nil {
		log.Printf("%d|%s() %v", user.ChatID, callingFunctionName, err)
	}
}

func DeleteLastMessage(user *objects.User) {
	botMessage := Database.GetRecordBotMessage(user)

	// Определение конфигурации для удаления сообщения
	deleteConfig := tapioca.DeleteMessageConfig{
		ChatID:    user.ChatID,
		MessageID: botMessage.MessageId, // Здесь нужно указать ID сообщения, которое вы хотите удалить
	}
	// Отправьте запрос на удаление сообщения
	_, err := Settings.GlobalBot.DeleteMessage(deleteConfig)
	if err != nil {
		log.Printf("%v", err)
	}

	Database.DeleteBotMessage(user)
}

// ProcessingMessage обрабатывает ответ пользователя на текстовое сообщения пользователя
func ProcessingMessage(update *tapioca.Update) {
	user, err := Database.GetUserFromDB(update)

	if err != nil {
		log.Printf("%s() %v", "ProcessingMessage", err)
	}

	if IsAddNewActivity(user) {
		// Пользователь нажал на кнопку добавления новой активности
		coins, text, err := scope.SplitSentence(update.Message.Text)

		if coins > 100 && coins > 0 {
			SendMessageWithDelay(user, 1, "Количество монет должно быть больше 0 и меньше 100")
		}

		if err != "" {
			SendMessageWithDelay(user, 1, err)
		} else {

			SendBotIsTyping(user)
			Delay(1)

			err := Database.AddNewActivity(user, coins, text)
			if err != nil {
				return
			} else {
				//Database.UpdateUserActionIdFromDB(user, update)
				ShowActivities(
					user,
					Database.GetActivitiesByUser(user),
					"Отследите полезную активность или вернитесь в главное меню:",
					update,
					true)
			}
		}
	} else if IsAddNewReward(user) {
		// Пользователь нажал на кнопку добавления новой награды
		coins, text, err := scope.SplitSentence(update.Message.Text)

		if coins > 100 && coins > 0 {
			SendMessageWithDelay(user, 1, "Количество монет должно быть больше 0 и меньше 100")
		}

		if err != "" {
			SendMessageWithDelay(user, 1, err)
		} else {

			SendBotIsTyping(user)
			Delay(1)

			err := Database.AddNewReward(user, coins, text)
			if err != nil {
				return
			} else {
				//Database.UpdateUserActionIdFromDB(user, update)
				ShowRewards(
					user,
					Database.GetRewardsByUser(user),
					"Отследите награду или вернитесь в главное меню:",
					update,
					true)
			}
		}
	}
}

// UpdateProcessing обрабатывает ответ пользователя на inline-кнопки
func UpdateProcessing(update *tapioca.Update) {
	// Структура User
	user := StoreUserFromUpdate(update)
	// Сохраните нового пользователя в базу данных
	err := Database.StoreUserToDB(user)
	if err != nil {
		log.Printf("%d|||||ERROR in UpdateProcessing(). %s. Ошибка записи данных о пользователе ", user.ChatID, err)
	}

	choiceCode := update.CallbackQuery.Data

	switch choiceCode {
	case Settings.BUTTON_CODE_BALANCE:
		DeleteLastMessage(user)
		ShowBalance(user, update)
	case Settings.BUTTON_CODE_ACTIVITIES: //Показать полезные активности
		DeleteLastMessage(user)
		ShowUsefulActivities(user, update, true)
	case Settings.BUTTON_CODE_REWARDS:
		DeleteLastMessage(user)
		ShowUsefulRewards(user, update, true)
	case Settings.BUTTON_CODE_PRINT_INTRO:
		PrintIntro(user)
		ShowMenu(user, update, true)
	case Settings.BUTTON_CODE_SKIP_INTRO: // Пропустить стартовые сообщения
		DeleteLastMessage(user)
		ShowMenu(user, update, true)
	case Settings.BUTTON_CODE_PRINT_MENU:
		DeleteLastMessage(user)
		ShowMenu(user, update, true)
	case Settings.BUTTON_CODE_ADD_ACTIVITY:
		DeleteLastMessage(user)

		SendBotIsTyping(user)
		SendMessageWithDelay(user, 1, "Введите количество монет, пробел и название активности. Пример: 2 Отжимания")

		Database.UpdateUserActionIdFromDB(user, update)
	case Settings.BUTTON_CODE_DELETE_ACTIVITY:
		Database.UpdateUserActionIdFromDB(user, update)
		ShowUsefulActivitiesAndDeleteBtn(user, update)
	case Settings.BUTTON_CODE_ADD_REWARD:
		DeleteLastMessage(user)

		SendBotIsTyping(user)
		SendMessageWithDelay(user, 1, "Введите количество монет, пробел и название награды. Пример: 2 Съесть кусок пиццы")

		Database.UpdateUserActionIdFromDB(user, update)
	case Settings.BUTTON_CODE_DELETE_REWARD:
		Database.UpdateUserActionIdFromDB(user, update)
		ShowUsefulRewardAndDeleteBtn(user, update)
	default:
		if IsDeleteActivity(user) {
			// Пользователь нажал на кнопку удаления активности
			DeleteActivity(update, user)
			ShowUsefulActivitiesAndDeleteBtn(user, update)
			return
		} else if IsDeleteReward(user) {
			// Пользователь нажал на кнопку удаления награды
			DeleteReward(update, user)
			ShowUsefulRewardAndDeleteBtn(user, update)
			return
		} else {
			// Обработка обычной активности
			// TODO:Переделать на поиск сразу в БД. И удалить метод FindActivity!!!
			if usefulActivity, found := FindActivity(Database.GetActivitiesByUser(user), choiceCode); found {
				ProcessUsefulActivity(usefulActivity, user)
				DeleteLastMessage(user)

				SendBotIsTyping(user)
				Delay(1)

				ShowUsefulActivities(user, update, true)
				return
			}

			// TODO:Переделать на поиск сразу в БД. И удалить метод FindReward!!!
			if reward, found := FindReward(Database.GetRewardsByUser(user), choiceCode); found {
				//TODO:Переделать ProcessUsefulReward или избавиться, т.к. работа идет с сущностью наград, а не с активностями
				ProcessUsefulReward(reward, user)
				DeleteLastMessage(user)

				SendBotIsTyping(user)
				Delay(1)

				ShowUsefulRewards(user, update, true)
				return
			}
		}

		log.Printf(`%d|||||[%T] !!!!!!!!! ERROR: Unknown code "%s"`, user.ChatID, time.Now(), choiceCode)
		msg := fmt.Sprintf("%s, Извините, я не распознаю код '%s' %s. Пожалуйста, сообщите об этой ошибке моему создателю.", user.FirstName+" "+user.LastName, choiceCode, Settings.EMOJI_SAD)
		SendStringMessage(user, msg)
	}
}

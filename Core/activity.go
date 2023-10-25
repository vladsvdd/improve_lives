package Core

import (
	"fmt"
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	//tapioca "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"improve_lives/Database"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
)

func DeleteActivity(update *tapioca.Update, user *objects.User) {
	choiceCode := update.CallbackQuery.Data
	// Создайте экземпляр структуры для обновления
	updateData := objects.Activity{
		IsDeleted: true,
	}
	// Обновите запись
	err := Settings.GlobalDatabase.Model(&objects.Activity{}).
		Where("user_id = ? AND code = ?", user.Id, choiceCode).
		Updates(updateData).Error

	if err != nil {
		log.Printf("%v", err)
	}
}

// ShowUsefulActivitiesAndDeleteBtn отправляет список полезных активностей для формы удаления
func ShowUsefulActivitiesAndDeleteBtn(user *objects.User, update *tapioca.Update) {
	ShowActivitiesAndDeleteBtn(user, Database.GetActivitiesByUser(user), "Удалите выбранную активность", update)
}

// ShowActivitiesAndDeleteBtn отправляет список активностей и кнопка удаления
func ShowActivitiesAndDeleteBtn(user *objects.User, activities []objects.Activity, message string, update *tapioca.Update) {
	activitiesButtonsRows := make([][]tapioca.InlineKeyboardButton, 0, len(activities)+1)
	for _, activity := range activities {
		activityDescription := ""
		activityDescription = fmt.Sprintf("%s + %d %s: %s", Settings.EMOJI_DELETE, activity.Coins, Settings.EMOJI_COIN, activity.Name)
		activitiesButtonsRows = append(activitiesButtonsRows, GetKeyboardRow(activityDescription, activity.Code))
	}
	activitiesButtonsRows = append(activitiesButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_ADD_ACTIVITY, Settings.BUTTON_CODE_ADD_ACTIVITY))
	activitiesButtonsRows = append(activitiesButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_DELETE_ACTIVITY, Settings.BUTTON_CODE_DELETE_ACTIVITY))
	activitiesButtonsRows = append(activitiesButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_PRINT_MENU, Settings.BUTTON_CODE_PRINT_MENU))

	buttons := tapioca.NewInlineKeyboardMarkup(activitiesButtonsRows...)
	SendUpdatedMessage(user, update, buttons, message, "ShowActivitiesAndDeleteBtn")
}

// GetUserActivityObj Получаем структуру действия пользователя
func GetUserActivityObj(user *objects.User) objects.UserAction {
	userAction := objects.UserAction{}
	err := Settings.GlobalDatabase.Table(objects.UserAction{}.TableName()).Where("id = ?", user.UserActionId).First(&userAction).Error
	if err != nil {
		log.Printf("%s() %v", "IsDeleteActivity", err)
	}

	return userAction
}

// IsDeleteActivity пользователь выбирал действие удаления?
func IsDeleteActivity(user *objects.User) bool {
	userAction := GetUserActivityObj(user)
	if Settings.BUTTON_CODE_DELETE_ACTIVITY == userAction.Action {
		return true
	}

	return false
}

// IsAddNewActivity пользователь выбирал действие добавления новой активности?
func IsAddNewActivity(user *objects.User) bool {
	userAction := GetUserActivityObj(user)

	if Settings.BUTTON_CODE_ADD_ACTIVITY == userAction.Action {
		return true
	}

	return false
}

func HandleActivityResult(user *objects.User, activity *objects.Activity, errorMsg string) {
	if activity.Coins == 0 {
		errorMsg = fmt.Sprintf(`действие "%s" не имеет указанной стоимости`, activity.Name)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, Извините, но %s %s Ваш баланс остается неизменным.", user.FirstName+" "+user.LastName, errorMsg, Settings.EMOJI_SAD)
	} else {
		err := Database.StoreActivityToDB(activity, user)
		if err != nil {
			log.Printf("%d|||||Ошибка записи активности в БД: %v", user.ChatID, err)
		} else {
			sumCoin, err := Database.GetSumCoinsActivityFromDB(user)
			if err != nil {
				log.Printf("%s() %v", "HandleActivityResult", err)
			}

			resultMessage = fmt.Sprintf(`%s, активность "%s" выполнена! %d %s были добавлены на ваш счет. Продолжайте в том же духе! %s%s Теперь у вас есть %d %s`,
				user.FirstName+" "+user.LastName,
				activity.Name,
				activity.Coins,
				Settings.EMOJI_COIN,
				Settings.EMOJI_BICEPS,
				Settings.EMOJI_SUNGLASSES,
				sumCoin,
				Settings.EMOJI_COIN)
		}
	}
	SendStringMessage(user, resultMessage)
}

// ShowActivities отправляет список активностей или наград пользователю
func ShowActivities(user *objects.User, activities []objects.Activity, message string, update *tapioca.Update, generateNewMessage bool) {
	activitiesButtonsRows := make([][]tapioca.InlineKeyboardButton, 0, len(activities)+1)
	for _, activity := range activities {
		activityDescription := ""
		activityDescription = fmt.Sprintf("+ %d %s: %s", activity.Coins, Settings.EMOJI_COIN, activity.Name)
		activitiesButtonsRows = append(activitiesButtonsRows, GetKeyboardRow(activityDescription, activity.Code))
	}

	activitiesButtonsRows = append(activitiesButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_ADD_ACTIVITY, Settings.BUTTON_CODE_ADD_ACTIVITY))
	activitiesButtonsRows = append(activitiesButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_DELETE_ACTIVITY, Settings.BUTTON_CODE_DELETE_ACTIVITY))
	activitiesButtonsRows = append(activitiesButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_PRINT_MENU, Settings.BUTTON_CODE_PRINT_MENU))

	if generateNewMessage {
		msg := tapioca.NewMessage(user.ChatID, message)
		msg.ReplyMarkup = tapioca.NewInlineKeyboardMarkup(activitiesButtonsRows...)
		message, err := Settings.GlobalBot.Send(msg)
		if err != nil {
			log.Printf("%d|||||ERROR in ShowActivities() %s", user.ChatID, err)
		}

		Database.CreateBotMessage(user, &message)
	} else {
		buttons := tapioca.NewInlineKeyboardMarkup(activitiesButtonsRows...)
		SendUpdatedMessage(user, update, buttons, message, "ShowUsefulActivities")
	}
}

// ShowUsefulActivities отправляет список полезных активностей
func ShowUsefulActivities(user *objects.User, update *tapioca.Update, generateNewMessage bool) {
	Database.UpdateUserActionIdFromDB(user, update)
	ShowActivities(
		user,
		Database.GetActivitiesByUser(user),
		"Отследите полезную активность или вернитесь в главное меню:",
		update,
		generateNewMessage)
}

// FindActivity ищет активность или награду по коду
func FindActivity(activities []objects.Activity, choiceCode string) (activity *objects.Activity, found bool) {
	for _, activity := range activities {
		if choiceCode == activity.Code {
			return &activity, true
		}
	}
	return
}

// ProcessUsefulActivity обрабатывает полезную активность
func ProcessUsefulActivity(activity *objects.Activity, user *objects.User) {
	errorMsg := ""
	if activity.Coins == 0 {
		errorMsg = fmt.Sprintf(`действие "%s" не имеет указанной стоимости`, activity.Name)
	}

	HandleActivityResult(user, activity, errorMsg)
}

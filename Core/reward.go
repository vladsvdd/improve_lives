package Core

import (
	"fmt"
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	"math"

	//tapioca "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"improve_lives/Database"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
)

func DeleteReward(update *tapioca.Update, user *objects.User) {
	choiceCode := update.CallbackQuery.Data
	// Создайте экземпляр структуры для обновления
	updateData := objects.Reward{
		IsDeleted: true,
	}
	// Обновите запись
	err := Settings.GlobalDatabase.Model(&objects.Reward{}).
		Where("user_id = ? AND code = ?", user.Id, choiceCode).
		Updates(updateData).Error

	if err != nil {
		log.Printf("%v", err)
	}
}

// ShowUsefulRewardAndDeleteBtn отправляет список наград для формы удаления
func ShowUsefulRewardAndDeleteBtn(user *objects.User, update *tapioca.Update) {
	ShowRewardAndDeleteBtn(user, Database.GetRewardsByUser(user), "Удалите выбранную активность", update)
}

// ShowRewardAndDeleteBtn отправляет список наград и кнопка удаления
func ShowRewardAndDeleteBtn(user *objects.User, activities []objects.Reward, message string, update *tapioca.Update) {
	rewardsButtonsRows := make([][]tapioca.InlineKeyboardButton, 0, len(activities)+1)
	for _, activity := range activities {
		rewardDescription := ""
		rewardDescription = fmt.Sprintf("%s + %d %s: %s", Settings.EMOJI_DELETE, activity.Coins, Settings.EMOJI_COIN, activity.Name)
		rewardsButtonsRows = append(rewardsButtonsRows, GetKeyboardRow(rewardDescription, activity.Code))
	}
	rewardsButtonsRows = append(rewardsButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_ADD_ACTIVITY, Settings.BUTTON_CODE_ADD_ACTIVITY))
	rewardsButtonsRows = append(rewardsButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_DELETE_ACTIVITY, Settings.BUTTON_CODE_DELETE_ACTIVITY))
	rewardsButtonsRows = append(rewardsButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_PRINT_MENU, Settings.BUTTON_CODE_PRINT_MENU))

	buttons := tapioca.NewInlineKeyboardMarkup(rewardsButtonsRows...)
	SendUpdatedMessage(user, update, buttons, message, "ShowRewardAndDeleteBtn")
}

// IsDeleteReward пользователь выбирал действие удаления?
func IsDeleteReward(user *objects.User) bool {
	userAction := GetUserActivityObj(user)
	if Settings.BUTTON_CODE_DELETE_REWARD == userAction.Action {
		return true
	}

	return false
}

// IsAddNewReward пользователь выбирал действие добавления новой награды?
func IsAddNewReward(user *objects.User) bool {
	userAction := GetUserActivityObj(user)

	if Settings.BUTTON_CODE_ADD_REWARD == userAction.Action {
		return true
	}

	return false
}

// ShowRewards TODO:Метод такой же как ShowActivities переделать под единый. Только работа с разными сущностями
// ShowRewards отправляет список активностей или наград пользователю
func ShowRewards(user *objects.User, rewards []objects.Reward, message string, update *tapioca.Update, generateNewMessage bool) {
	rewardsButtonsRows := make([][]tapioca.InlineKeyboardButton, 0, len(rewards)+1)
	for _, reward := range rewards {
		rewardDescription := ""
		rewardDescription = fmt.Sprintf("- %d %s: %s", reward.Coins, Settings.EMOJI_COIN, reward.Name)
		rewardsButtonsRows = append(rewardsButtonsRows, GetKeyboardRow(rewardDescription, reward.Code))
	}
	rewardsButtonsRows = append(rewardsButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_ADD_REWARD, Settings.BUTTON_CODE_ADD_REWARD))
	rewardsButtonsRows = append(rewardsButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_DELETE_REWARD, Settings.BUTTON_CODE_DELETE_REWARD))
	rewardsButtonsRows = append(rewardsButtonsRows, GetKeyboardRow(Settings.BUTTON_TEXT_PRINT_MENU, Settings.BUTTON_CODE_PRINT_MENU))

	if generateNewMessage {
		msg := tapioca.NewMessage(user.ChatID, message)
		msg.ReplyMarkup = tapioca.NewInlineKeyboardMarkup(rewardsButtonsRows...)
		message, err := Settings.GlobalBot.Send(msg)
		if err != nil {
			log.Printf("%d|||||ERROR in ShowRewards() %s", user.ChatID, err)
		}

		Database.CreateBotMessage(user, &message)
	} else {
		buttons := tapioca.NewInlineKeyboardMarkup(rewardsButtonsRows...)
		SendUpdatedMessage(user, update, buttons, message, "ShowRewards")
	}
}

// ShowUsefulRewards отправляет список наград
func ShowUsefulRewards(user *objects.User, update *tapioca.Update, generateNewMessage bool) {
	Database.UpdateUserActionIdFromDB(user, update)
	ShowRewards(
		user,
		Database.GetRewardsByUser(user),
		"Отследите награду или вернитесь в главное меню:",
		update,
		generateNewMessage)
}

// FindReward ищет активность или награду по коду
func FindReward(rewards []objects.Reward, choiceCode string) (reward *objects.Reward, found bool) {
	for _, reward := range rewards {
		if choiceCode == reward.Code {
			return &reward, true
		}
	}
	return
}

func HandleRewardResult(user *objects.User, reward *objects.Reward, errorMsg string) {
	if reward.Coins == 0 {
		errorMsg = fmt.Sprintf(`действие "%s" не имеет указанной стоимости`, reward.Name)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, Извините, но %s %s Ваш баланс остается неизменным.", user.FirstName+" "+user.LastName, errorMsg, Settings.EMOJI_SAD)
	} else {
		err := Database.StoreRewardToDB(reward, user)
		if err != nil {
			log.Printf("%d|||||Ошибка записи активности в БД: %v", user.ChatID, err)
		} else {
			sumCoin, err := Database.GetSumCoinsActivityFromDB(user)
			if err != nil {
				log.Printf("%s() %v", "HandleActivityResult", err)
			}

			resultMessage = fmt.Sprintf(`%s, Награда "%s" оплачена, начинайте! С вашего счета списано %d %s. Теперь у вас есть %d %s`,
				user.FirstName+" "+user.LastName,
				reward.Name,
				int(math.Abs(float64(reward.Coins))),
				Settings.EMOJI_COIN,
				sumCoin,
				Settings.EMOJI_COIN)
		}
	}
	SendStringMessage(user, resultMessage)
}

// ProcessUsefulReward обрабатывает награду
func ProcessUsefulReward(reward *objects.Reward, user *objects.User) {
	errorMsg := ""
	sumCoin, err := Database.GetSumCoinsActivityFromDB(user)
	if err != nil {
		log.Printf("%s() %v", "ProcessUsefulReward", err)
	}
	if reward.Coins == 0 {
		errorMsg = fmt.Sprintf(`Награда "%s" не имеет указанной стоимости.`, reward.Name)
	} else if sumCoin < reward.Coins {
		errorMsg = fmt.Sprintf(`В данный момент у вас есть %d %s. Вы не можете позволить себе "%s" за %d %s.`, sumCoin, Settings.EMOJI_COIN, reward.Name, reward.Coins, Settings.EMOJI_COIN)
	}

	HandleRewardResult(user, reward, errorMsg)
}

package Core

import (
	"fmt"
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	"improve_lives/Database"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
)

// PrintIntro отправляет вступительные сообщения
func PrintIntro(user *objects.User) {
	SendMessageWithDelay(user, 2, "Привет! "+Settings.EMOJI_SUNGLASSES)
	SendMessageWithDelay(user, 7, "Существует множество полезных действий, которые, выполняя их регулярно, улучшают качество нашей жизни. Однако зачастую бывает более весело, легко или вкусно делать что-то вредное. Не так ли?")
	SendMessageWithDelay(user, 7, "С большей вероятностью мы предпочтем заблудиться в YouTube Shorts вместо урока английского, купить M&M's вместо овощей или лежать в постели вместо занятий йогой.")
	SendMessageWithDelay(user, 1, Settings.EMOJI_SAD)
	SendMessageWithDelay(user, 10, "Каждый играл хотя бы в одну игру, где нужно повышать уровень персонажа, делая его сильнее, умнее или красивее. Это приносит удовольствие, потому что каждое действие приносит результаты. В реальной жизни систематические действия со временем начинают замечаться. Давайте это изменить, не так ли?")
	SendMessageWithDelay(user, 1, Settings.EMOJI_SMILE)
	SendMessageWithDelay(user, 14, `Перед вами две таблицы: "Полезные действия" и "Награды". Первая таблица содержит список коротких простых действий, и за выполнение каждого из них вы получите указанное количество монет. Во второй таблице вы увидите список действий, которые можно выполнить только после оплаты их монетами, заработанными на предыдущем этапе.`)
	SendMessageWithDelay(user, 1, Settings.EMOJI_COIN)
	SendMessageWithDelay(user, 10, `Например, вы тратите полчаса на занятия йогой, за что получаете 2 монеты. Затем у вас есть 2 часа занятий программированием, за что вы получаете 8 монет. Теперь вы можете посмотреть 1 эпизод "Интернов" и остаться при этом в нуле. Просто и понятно!`)
	SendMessageWithDelay(user, 6, `Помечайте выполненные полезные действия, чтобы не потерять свои монеты. И не забудьте "купить" награду, прежде чем ее выполнить на самом деле.`)
}

// AskToPrintIntro отправляет сообщение с вопросом о прочтении вступительных сообщений
func AskToPrintIntro(user *objects.User) {
	msg := tapioca.NewMessage(user.ChatID, "Хочешь прочитать вступительные сообщения и узнать правила игры?")
	msg.ReplyMarkup = tapioca.NewInlineKeyboardMarkup(
		GetKeyboardRow(Settings.BUTTON_TEXT_PRINT_INTRO, Settings.BUTTON_CODE_PRINT_INTRO),
		GetKeyboardRow(Settings.BUTTON_TEXT_SKIP_INTRO, Settings.BUTTON_CODE_SKIP_INTRO),
	)
	message, err := Settings.GlobalBot.Send(msg)
	if err != nil {
		log.Printf("%d||||| ERROR in %s() %v", user.ChatID, "AskToPrintIntro", err)
	}

	Database.CreateBotMessage(user, &message)
}

// ShowMenu отправляет пользователю меню с кнопками
func ShowMenu(user *objects.User, update *tapioca.Update, generateNewMessage bool) {
	message := "Выберите один из вариантов:"
	buttons := tapioca.NewInlineKeyboardMarkup(
		GetKeyboardRow(Settings.BUTTON_TEXT_BALANCE, Settings.BUTTON_CODE_BALANCE),
		GetKeyboardRow(Settings.BUTTON_TEXT_ACTIVITIES, Settings.BUTTON_CODE_ACTIVITIES),
		GetKeyboardRow(Settings.BUTTON_TEXT_REWARDS, Settings.BUTTON_CODE_REWARDS),
	)
	if generateNewMessage {
		msg := tapioca.NewMessage(user.ChatID, message)
		msg.ReplyMarkup = buttons
		message, err := Settings.GlobalBot.Send(msg)

		if err != nil {
			log.Printf("%d||||| ERROR in %s() %v", user.ChatID, "AskToPrintIntro", err)
		}

		Database.CreateBotMessage(user, &message)
	} else {
		SendUpdatedMessage(user, update, buttons, message, "ShowMenu")
	}
}

// ShowBalance отправляет баланс пользователя
func ShowBalance(user *objects.User, update *tapioca.Update) {
	msg := fmt.Sprintf("%s, в данный момент ваш кошелек пуст %s \nОтслеживайте полезные действия, чтобы зарабатывать монеты", user.FirstName+" "+user.LastName, Settings.EMOJI_DONT_KNOW)
	coins, _ := Database.GetSumCoinsActivityFromDB(user)
	if coins > 0 {
		msg = fmt.Sprintf("%s, у тебя %d %s", user.FirstName+" "+user.LastName, coins, Settings.EMOJI_COIN)
	}
	SendStringMessage(user, msg)
	ShowMenu(user, update, true)
}

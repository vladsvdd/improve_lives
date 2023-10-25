package Settings

const (
	EMOJI_THUMBS_UP    = "\U0001F44D"   // 👍
	EMOJI_SMILE        = "\U0001F604"   // 😄
	EMOJI_HEART        = "\U00002764"   // ❤️
	EMOJI_ROCKET       = "\U0001F680"   // 🚀
	EMOJI_DELETE       = "\U0000274C"   // ❌
	EMOJI_CAMERA       = "\U0001F4F7"   // 📷
	EMOJI_BOOK         = "\U0001F4D6"   // 📖
	EMOJI_MUSIC_NOTE   = "\U0001F3B5"   // 🎵
	EMOJI_WINK         = "\U0001F609"   // 😉
	EMOJI_CHEESEBURGER = "\U0001F354"   // 🍔
	EMOJI_SUN          = "\U00002600"   // ☀️
	EMOJI_LEFT_ARROW   = "\U00002B05"   // ⬅️
	EMOJI_RIGHT_ARROW  = "\U000027A1"   // ➡️
	EMOJI_COIN         = "\U0001FA99"   // (coin)
	EMOJI_DOLLAR       = "\U0001F4B5"   // 💵
	EMOJI_SUNGLASSES   = "\U0001F60E"   // 😎
	EMOJI_WOW          = "\U0001F604"   // 😄
	EMOJI_DONT_KNOW    = "\U0001F937"   // 🤷
	EMOJI_SAD          = "\U0001F63F"   // 😿
	EMOJI_BICEPS       = "\U0001F4AA"   // 💪
	EMOJI_BUTTON_START = "\U000025B6  " // ▶
	EMOJI_BUTTON_END   = "  \U000025C0" // ◀

	TOKEN_NAME_IN_OS            = "TOKEN_NAME_IN_OS"
	DB_USERNAME                 = "DB_USERNAME"
	DB_PASSWORD                 = "DB_PASSWORD"
	DB_HOST                     = "DB_HOST"
	DB_PORT                     = "DB_PORT"
	DB_DATABASE                 = "DB_DATABASE"
	UPDATE_CONFIG_TIMEOUT       = 60
	MAX_USER_COINS        int16 = 500

	BUTTON_TEXT_PRINT_INTRO     = EMOJI_BUTTON_START + "Просмотреть введение" + EMOJI_BUTTON_END
	BUTTON_TEXT_SKIP_INTRO      = EMOJI_BUTTON_START + "Пропустить введение" + EMOJI_BUTTON_END
	BUTTON_TEXT_BALANCE         = EMOJI_BUTTON_START + "Текущий баланс" + EMOJI_BUTTON_END
	BUTTON_TEXT_ACTIVITIES      = EMOJI_BUTTON_START + "Полезные действия" + EMOJI_BUTTON_END
	BUTTON_TEXT_REWARDS         = EMOJI_BUTTON_START + "Награды" + EMOJI_BUTTON_END
	BUTTON_TEXT_PRINT_MENU      = EMOJI_BUTTON_START + "ГЛАВНОЕ МЕНЮ" + EMOJI_BUTTON_END
	BUTTON_TEXT_ADD_ACTIVITY    = EMOJI_BUTTON_START + "ДОБАВИТЬ СВОЮ АКТИВНОСТЬ" + EMOJI_BUTTON_END
	BUTTON_TEXT_DELETE_ACTIVITY = EMOJI_BUTTON_START + "УДАЛИТЬ АКТИВНОСТЬ" + EMOJI_BUTTON_END
	BUTTON_TEXT_ADD_REWARD      = EMOJI_BUTTON_START + "ДОБАВИТЬ СВОЮ НАГРАДУ" + EMOJI_BUTTON_END
	BUTTON_TEXT_DELETE_REWARD   = EMOJI_BUTTON_START + "УДАЛИТЬ НАГРАДУ" + EMOJI_BUTTON_END

	BUTTON_CODE_PRINT_INTRO = "print_intro"
	BUTTON_CODE_SKIP_INTRO  = "skip_intro"
	BUTTON_CODE_BALANCE     = "show_balance"
	BUTTON_CODE_PRINT_MENU  = "print_menu"

	BUTTON_CODE_ACTIVITIES      = "show_activities"
	BUTTON_CODE_ADD_ACTIVITY    = "add_activity"
	BUTTON_CODE_DELETE_ACTIVITY = "delete_activity"

	BUTTON_CODE_REWARDS       = "show_rewards"
	BUTTON_CODE_ADD_REWARD    = "add_reward"
	BUTTON_CODE_DELETE_REWARD = "delete_reward"
)

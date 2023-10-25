package scope

import (
	"strconv"
	"strings"
)

func SplitSentence(input string) (int, string, string) {
	// Разделите строку на слова
	words := strings.Fields(input)
	// Проверьте, есть ли хотя бы одно слово
	if len(words) > 0 {
		// Попробуйте извлечь число из первого слова
		if num, err := strconv.Atoi(words[0]); err == nil {
			// Извлеките число из первого слова
			message := strings.Join(words[1:], " ") // Объедините остальные слова обратно в сообщение

			// Выведите число и сообщение
			return num, message, ""
		} else {
			return 0, "", "Сообщение без числа"
		}
	}

	return 0, "", "Сообщение не может быть пустым"
}

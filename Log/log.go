package Log

import (
	"log"
	"os"
)

func InitErrorFile() *os.File {
	// Открываем файл для записи логов. Если файл не существует, он будет создан.
	logFile, err := os.OpenFile("./Log/errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Устанавливаем файл как вывод для логов
	log.SetOutput(logFile)

	// Не забудьте закрыть файл после использования
	return logFile
}

package Settings

import (
	"github.com/joho/godotenv"
	"log"
)

// LoadEnv Загрузите переменные среды из файла .env
func LoadEnv(filename string) {
	if err := godotenv.Load(filename); err != nil {
		log.Fatalf("Ошибка при загрузке файла .env: %s", err)
	}
}

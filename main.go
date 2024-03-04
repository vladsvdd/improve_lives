package main

import (
	"github.com/gin-gonic/gin"
	tapioca "github.com/go-telegram-bot-api/telegram-bot-api"
	//tapioca "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"improve_lives/Core"
	"improve_lives/Database"
	"improve_lives/Log"
	"improve_lives/Settings"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	// Инициализируйте мьютекс для безопасного доступа к данным.
	_ sync.Mutex
)

func processUpdates(updates <-chan tapioca.Update) {
	for update := range updates {
		if Core.IsCallbackQuery(&update) {
			// Обработайте callback-запрос в этой функции.
			go Core.UpdateProcessing(&update)
		} else if Core.IsStartMessage(&update) {
			// Обработайте стартовое сообщение в этой функции.
			go startMessageProcessing(&update)
		} else if update.Message != nil {
			// Обработка текстовых сообщений
			go Core.ProcessingMessage(&update)
		}
	}
}

func startMessageProcessing(update *tapioca.Update) {
	// Здесь вы можете проверить, есть ли пользователь в базе данных и загрузить его данные
	user, err := Database.GetUserFromDB(update)

	if err != nil {
		log.Printf("%s() %v", "startMessageProcessing", err)
	}

	Core.AskToPrintIntro(user)
}

func initEnv() {
	var err error

	//.env
	if Settings.GlobalTest {
		Settings.LoadEnv("test.env")
	} else {
		Settings.LoadEnv(".env")
	}

	Settings.GlobalPort = os.Getenv("PORT")
	if len(Settings.GlobalPort) == 0 {
		log.Panicf("%d|||||%s() Не удалось загрузить переменную среды PORT", 0, "init")
	}

	// Получите токен из переменной среды
	if Settings.GlobalToken = os.Getenv(Settings.TOKEN_NAME_IN_OS); Settings.GlobalToken == "" {
		log.Panicf("%d|||||%s() Не удалось загрузить переменную среды Settings.TOKEN_NAME_IN_OS", 0, "init")
	}

	if Settings.GlobalBot, err = tapioca.NewBotAPI(Settings.GlobalToken); err != nil {
		log.Panicf("%d|||||%s() %s", 0, "init", err)
	}
}

func init() {
	//Settings.GlobalTest меняем на false перед выкладкой в прод
	Settings.GlobalTest = true

	//.env
	initEnv()

	var err error
	// Открываем соединение с базой данных SQLite
	db, err := Database.InitMySQLDB()
	if err != nil {
		log.Fatalf("%d|||||%s ERROR in init()", 0, err)
	}
	Settings.GlobalDatabase = db

}

func main() {
	// Инициализируйте логгер ошибок.
	logFile := Log.InitErrorFile()
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Printf("%d|||||ERROR in main 'defer func(logFile *os.File)' %s", 0, err)
		}
	}(logFile) // Закрываем лог файл

	// Создайте экземпляр маршрутизатора Gin и определите маршрут.
	router := gin.Default()
	router.GET("/improve_lives/", func(c *gin.Context) {
		c.String(http.StatusOK, "Привет, мир!")
	})

	// Инициализируйте канал для обработки обновлений.
	updates := make(chan tapioca.Update)
	go processUpdates(updates)

	// Запустите HTTP-сервер и обработку обновлений в отдельных горутинах.
	go func() {
		if err := http.ListenAndServe(":"+Settings.GlobalPort, router); err != nil {
			log.Printf("%d|||||%s() %s", 0, "main", err)
		}
		err := router.Run()
		if err != nil {
			log.Printf("%d|||||%s() -> r.Run(). %s", 0, "main", err)
		}
	}()

	// Начните получение обновлений от бота.
	updateConfig := tapioca.NewUpdate(0)
	updateConfig.Timeout = Settings.UPDATE_CONFIG_TIMEOUT

	updatesChan, _ := Settings.GlobalBot.GetUpdatesChan(updateConfig)

	for update := range updatesChan {
		// Отправьте обновление в канал для обработки.
		updates <- update
	}
}

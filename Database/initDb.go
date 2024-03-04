package Database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"improve_lives/Settings"
	"improve_lives/db/objects"
	"log"
	"os"
)

func InitMySQLDB() (*gorm.DB, error) {
	// Получите токен из переменной среды
	if Settings.GlobalDbUsername = os.Getenv(Settings.DB_USERNAME); Settings.GlobalDbUsername == "" {
		log.Panicf("%d|||||%s() Не удалось загрузить переменную среды Settings.DB_USERNAME", 0, "InitMySQLDB")
	}
	if Settings.GlobalDbPassword = os.Getenv(Settings.DB_PASSWORD); Settings.GlobalDbPassword == "" {
		log.Panicf("%d|||||%s() Не удалось загрузить переменную среды Settings.GlobalDB_PASSWORD", 0, "InitMySQLDB")
	}
	if Settings.GlobalDbHost = os.Getenv(Settings.DB_HOST); Settings.GlobalDbHost == "" {
		log.Panicf("%d|||||%s() Не удалось загрузить переменную среды Settings.DB_HOST", 0, "InitMySQLDB")
	}
	if Settings.GlobalDbPort = os.Getenv(Settings.DB_PORT); Settings.GlobalDbPort == "" {
		log.Panicf("%d|||||%s() Не удалось загрузить переменную среды Settings.DB_PORT", 0, "InitMySQLDB")
	}
	if Settings.GlobalDbDatabase = os.Getenv(Settings.DB_DATABASE); Settings.GlobalDbDatabase == "" {
		log.Panicf("%d|||||%s() Не удалось загрузить переменную среды Settings.DB_DATABASE", 0, "InitMySQLDB")
	}

	// Создайте строку подключения к базе данных MySQL.
	dbString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		Settings.GlobalDbUsername,
		Settings.GlobalDbPassword,
		//Settings.GlobalDbHost,
		"localhost",
		Settings.GlobalDbPort,
		Settings.GlobalDbDatabase)
	fmt.Println(dbString)
	// Откройте соединение с базой данных
	db, err := gorm.Open(mysql.Open(dbString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

// InsertActivities добавление активностей пользователя
func InsertActivities(activities *[]objects.Activity) error {
	// Начинаем транзакцию
	tx := Settings.GlobalDatabase.Begin()

	// Проверяем, была ли ошибка при начале транзакции
	if tx.Error != nil {
		log.Println(tx.Error)
	}

	// Добавляем записи внутри транзакции
	for _, activity := range *activities {
		if err := tx.Create(&activity).Error; err != nil {
			// Если произошла ошибка при вставке записи, откатываем транзакцию
			tx.Rollback()
			log.Println("Ошибка при вставке записи в базу данных:", err)
			return err
		}
	}

	// Если все успешно, фиксируем транзакцию
	tx.Commit()

	return nil
}

// InsertRewards добавление наград пользователя
func InsertRewards(rewards *[]objects.Reward) error {
	// Начинаем транзакцию
	tx := Settings.GlobalDatabase.Begin()

	// Проверяем, была ли ошибка при начале транзакции
	if tx.Error != nil {
		log.Println(tx.Error)
	}

	// Добавляем записи внутри транзакции
	for _, reward := range *rewards {
		if err := tx.Create(&reward).Error; err != nil {
			// Если произошла ошибка при вставке записи, откатываем транзакцию
			tx.Rollback()
			log.Println("Ошибка при вставке записи в базу данных:", err)
			return err
		}
	}

	// Если все успешно, фиксируем транзакцию
	tx.Commit()

	return nil
}

// InsertStartActivitiesAndRewards запись списка активностей и наград для нового пользователя
func InsertStartActivitiesAndRewards(user *objects.User) {
	if user.Id == 0 {
		log.Printf("ERROR in InsertStartActivitiesAndRewards(). Данные о пользователе: ID: %d не найдены в БД", user.Id)
		return
	}
	// Пример данных для вставки.
	activities := []objects.Activity{
		{UserId: user.Id, Code: "yoga", Name: "Йога (15 минут)", Coins: 5},
		{UserId: user.Id, Code: "meditation", Name: "Медитация (15 минут)", Coins: 5},
		{UserId: user.Id, Code: "language", Name: "Изучение иностранного языка (15 минут)", Coins: 4},
		{UserId: user.Id, Code: "swimming", Name: "Плавание (15 минут)", Coins: 5},
		{UserId: user.Id, Code: "walk", Name: "Прогулка (15 минут)", Coins: 2},
		{UserId: user.Id, Code: "chores", Name: "Домашние дела", Coins: 4},
		{UserId: user.Id, Code: "work_learning", Name: "Изучение рабочих материалов (15 минут)", Coins: 3},
		{UserId: user.Id, Code: "portfolio_work", Name: "Работа над проектом портфолио (15 минут)", Coins: 3},
		{UserId: user.Id, Code: "resume_edit", Name: "Редактирование резюме (15 минут)", Coins: 3},
		{UserId: user.Id, Code: "creative", Name: "Творческое творчество (15 минут)", Coins: 3},
		{UserId: user.Id, Code: "reading", Name: "Чтение художественной литературы (15 минут)", Coins: 3},
	}

	rewards := []objects.Reward{
		{UserId: user.Id, Code: "watch_series", Name: "Просмотр сериала (1 серию)", Coins: 10},
		{UserId: user.Id, Code: "watch_movie", Name: "Просмотр фильма (1 шт)", Coins: 30},
		{UserId: user.Id, Code: "social_nets", Name: "Просмотр социальных сетей (30 минут)", Coins: 10},
		{UserId: user.Id, Code: "eat_sweets", Name: "Съесть сладкое", Coins: 60},
	}

	// Вызываем метод для вставки данных.
	if len(GetActivitiesByUser(user)) == 0 {
		if err := InsertActivities(&activities); err != nil {
			log.Printf("%d|%v", user.Id, err)
		}
	}

	// Вызываем метод для вставки данных.
	if len(GetRewardsByUser(user)) == 0 {
		if err := InsertRewards(&rewards); err != nil {
			log.Printf("%d|%v", user.Id, err)
		}
	}
}

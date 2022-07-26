package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	pig_bot "pig-bot/pigs"
	"time"
)

const BotEnvVarName = "BOT_TOKEN"

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	dsn := fmt.Sprintf("host=database user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Moscow",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	pig_bot.AutoMigrate(db)
	botToken, ok := os.LookupEnv(BotEnvVarName)
	if !ok {
		panic(fmt.Sprintf("%s env variable must be set", BotEnvVarName))
	}

	bot, err := pig_bot.NewBot(&pig_bot.Params{Token: botToken, DB: db})
	bot.Start()
}

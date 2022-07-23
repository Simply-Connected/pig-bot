package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	pig_bot "pig-bot/pigs"
	"time"
	"fmt"
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
	dsn := "host=localhost user=postgres password=postgres dbname=tst port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	pig_bot.AutoMigrate(db)
	bot_token, ok := os.LookupEnv(BotEnvVarName)
	if !ok {
            panic(fmt.Sprintf("%s env variable must be set", BotEnvVarName))
	}


	bot, err := pig_bot.NewBot(&pig_bot.Params{Token: bot_token, DB: db})
	bot.Start()
}

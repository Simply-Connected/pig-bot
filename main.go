package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	pig_bot "pig-bot/pigs"
	"time"
)

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
	dsn := "host=localhost user=skazzi dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	pig_bot.AutoMigrate(db)
	bot, err := pig_bot.NewBot(&pig_bot.Params{Token: "5570481617:AAHi7W1XaFvwJSSGN_FnE8XJGk5QPv0NFkY", DB: db})
	bot.Start()
}

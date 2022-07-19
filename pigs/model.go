package pig_bot

import (
	"database/sql"
	"gorm.io/gorm"
)

type Pig struct {
	gorm.Model
	Weight   uint32
	Name     string
	LastGrow sql.NullTime
	UserID   int64
}

type User struct {
	gorm.Model
	TelegramId int64
	Name       string
	WinsCount  uint32
	Pig        Pig
}

type PigPicture struct {
	gorm.Model
	Data []byte
	Type ActionType
}

type ActionType int64

const (
	Untyped ActionType = iota
	Eating
	Fighting
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Pig{})
	db.AutoMigrate(&User{})
}

func GetUser(tgID int64, db *gorm.DB) (user *User, inserted bool) {
	user = &User{TelegramId: tgID}
	inserted = db.Where(user, "TelegramId").Preload("Pig").First(user).RowsAffected != 0
	return
}

func RegisterUser(tgID int64, tgName string, pigName string, db *gorm.DB) *User {
	if pigName == "" {
		pigName = "Свин"
	}
	user := &User{TelegramId: tgID,
		Name:      tgName,
		WinsCount: 0,
		Pig: Pig{
			Name:   pigName,
			Weight: 7,
		},
	}
	db.Create(user)
	return user
}

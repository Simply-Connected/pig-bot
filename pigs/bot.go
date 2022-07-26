package pig_bot

import (
	"database/sql"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Params struct {
	DB    *gorm.DB
	Token string
}

// /grow /rename /duel /top

func NewBot(params *Params) (*tele.Bot, error) {
	pref := tele.Settings{
		Token:  params.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	db := params.DB

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	b.Use(middleware.Logger())
	b.Use(middleware.AutoRespond())

	random := rand.New(rand.NewSource(time.Now().Unix()))

	b.Handle("/grow", func(c tele.Context) error {
		user := GetOrRegisterUser(c, db)

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return err
		}
		curTime := time.Now().In(loc)
		fmt.Printf("DEBUG:%s", curTime.UTC().String())
		if user.Pig.LastGrow.Valid &&
			curTime.Day() == user.Pig.LastGrow.Time.Day() &&
			curTime.Sub(user.Pig.LastGrow.Time).Hours() < 24 {
			return c.Send(fmt.Sprintf("Вы уже кормили свою свинью сегодня.\n\nВес вашего свина: *%d*", user.Pig.Weight), tele.ModeMarkdown)
		}

		diff := random.Int31n(16) - 5
		if int32(user.Pig.Weight)+diff <= 0 {
			user.Pig.Weight = 7
			user.Pig.LastGrow = sql.NullTime{Valid: false}
			db.Save(&user.Pig)
			return c.Send(fmt.Sprintf("Ваша свинья *%s* потеряла весь вес и умерла ☠️. "+
				"Вы получаете нового порося.\n\n"+
				"Вес вашего свина: *%d*",
				user.Pig.Name, user.Pig.Weight), tele.ModeMarkdown)
		}
		user.Pig.Weight = uint32(int32(user.Pig.Weight) + diff)
		user.Pig.LastGrow = sql.NullTime{Time: time.Now(), Valid: true}
		db.Save(&user.Pig)
		return c.Send(getGrowPhrase(&user.Pig, diff), tele.ModeMarkdown)
	})

	b.Handle("/rename", func(c tele.Context) error {
		if c.Message().Payload == "" {
			return c.Send("Вы ввели пустое имя. У свинки должно быть полноценное имя")
		}
		user := GetOrRegisterUser(c, db)

		user.Pig.Name = c.Message().Payload
		db.Save(&user.Pig)
		return c.Send(fmt.Sprintf("Теперь вашу свинью зовут *%s*", user.Pig.Name), tele.ModeMarkdown)
	})

	b.Handle("/top", func(c tele.Context) error {
		var pigs []Pig
		db.Order("weight desc").Limit(10).Find(&pigs)

		var message strings.Builder
		message.WriteString("Топ *10* свинок:\n\n")

		for i, pig := range pigs {
			message.WriteString(fmt.Sprintf("*%d*. *%s* - %d кг\n", i+1, pig.Name, pig.Weight))
		}
		return c.Send(message.String(), tele.ModeMarkdown)
	})
	return b, nil
}

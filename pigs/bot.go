package pig_bot

import (
	"database/sql"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"gorm.io/gorm"
	"log"
	"math/rand"
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
		tgID := c.Sender().ID
		user, inserted := GetUser(tgID, db)
		if !inserted {
			user = RegisterUser(tgID, c.Sender().Username, c.Sender().Username, db)
		}

		curTime := time.Now()
		if user.Pig.LastGrow.Valid &&
			curTime.Day() == user.Pig.LastGrow.Time.Day() &&
			curTime.Sub(user.Pig.LastGrow.Time).Hours() < 24 {
			c.Send("Вы уже кормили свою свинью сегодня.")
			return nil
		}

		diff := random.Int31n(16) - 5
		if int32(user.Pig.Weight)+diff <= 0 {
			user.Pig.Weight = 0
			user.Pig.LastGrow = sql.NullTime{Valid: false}
			db.Save(&user.Pig)
			c.Send(fmt.Sprintf("Ваша свинья %s потеряла весь вес и умерла ☠️. Попробуйте заново.", user.Pig.Name))
			return nil
		}
		user.Pig.Weight = uint32(int32(user.Pig.Weight) + diff)
		user.Pig.LastGrow = sql.NullTime{Time: time.Now(), Valid: true}
		db.Save(&user.Pig)
		if diff < 0 {
			c.Send(fmt.Sprintf("diff = %d, total = %d", diff, user.Pig.Weight))
		} else if diff == 0 {
			c.Send(fmt.Sprintf("diff = %d, total = %d", diff, user.Pig.Weight))
		} else {
			c.Send(fmt.Sprintf("diff = %d, total = %d", diff, user.Pig.Weight))
		}
		return nil
	})
	return b, nil
}

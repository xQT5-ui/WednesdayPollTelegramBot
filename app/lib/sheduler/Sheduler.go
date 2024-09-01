package sheduler

import (
	"time"

	conf "app.go/app/config"
	"app.go/app/lib"
	lg "app.go/app/lib/logger"
	"github.com/carlescere/scheduler"
	tb "gopkg.in/telebot.v3"
)

type Sheduler struct {
	Scheduler *scheduler.Job
}

func NewScheduler() *Sheduler {
	return &Sheduler{}
}

func (s *Sheduler) WednesdayJob(bot *tb.Bot, fact string, log *lg.Logger, config *conf.Config) *scheduler.Job {
	wednesdayJob, err := scheduler.Every().Day().At("12:20").Run(func() {
		// check day of week because sheduler.Every()DayOfWeek().Run() doesn't work
		if time.Now().Weekday() == time.Wednesday {
			lib.SendPoll(bot, &tb.Chat{ID: int64(config.Bot_secure.Chat_id)}, fact, log, config)
		}
	})
	if err != nil {
		log.Fatal(err, "Ошибка при планировании задачи на среду:")
	}

	log.Info("Запланирована задача на среду")

	return wednesdayJob
}

func (s *Sheduler) ThursdayJob(bot *tb.Bot, fact string, log *lg.Logger, config *conf.Config) *scheduler.Job {
	thursdayJob, err := scheduler.Every().Day().At("13:00").Run(func() {
		// check day of week because sheduler.Every()DayOfWeek().Run() doesn't work
		if time.Now().Weekday() == time.Thursday {
			lib.UnpinMsg(bot, &tb.Chat{ID: int64(config.Bot_secure.Chat_id)}, log, config)
		}
	})
	if err != nil {
		log.Fatal(err, "Ошибка при планировании задачи на четверг:")
	}

	log.Info("Запланирована задача на четверг")

	return thursdayJob
}

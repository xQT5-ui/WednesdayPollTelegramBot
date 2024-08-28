package sheduler

import (
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
	// FOR TEST
	// wednesdayJob, err := scheduler.Every(3).Seconds().NotImmediately().Run(func() {
	wednesdayJob, err := scheduler.Every().Wednesday().At("12:20").Run(func() {
		lib.SendPoll(bot, &tb.Chat{ID: int64(config.Bot_secure.Chat_id)}, fact, log, config)
	})
	if err != nil {
		log.Fatal(err, "Ошибка при планировании задачи на среду:")
	}

	log.Info("Запланирована задача на среду")

	return wednesdayJob
}

func (s *Sheduler) ThursdayJob(bot *tb.Bot, fact string, log *lg.Logger, config *conf.Config) *scheduler.Job {
	// FOR TEST
	// thursdayJob, err := scheduler.Every(4).Seconds().NotImmediately().Run(func() {
	thursdayJob, err := scheduler.Every().Thursday().At("13:00").Run(func() {
		lib.UnpinMsg(bot, &tb.Chat{ID: int64(config.Bot_secure.Chat_id)}, log, config)
	})
	if err != nil {
		log.Fatal(err, "Ошибка при планировании задачи на четверг:")
	}

	log.Info("Запланирована задача на четверг")

	return thursdayJob
}

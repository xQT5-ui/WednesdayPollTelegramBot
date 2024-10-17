package sheduler

import (
	"time"

	conf "app.go/app/config"
	lg "app.go/app/lib/logger"
	tgbot "app.go/app/lib/tgbot"

	"github.com/carlescere/scheduler"
	tb "gopkg.in/telebot.v3"
)

type Sheduler struct {
	Scheduler *scheduler.Job
}

func NewScheduler() *Sheduler {
	return &Sheduler{}
}

func (s *Sheduler) WednesdayJob(bot *tb.Bot, fact string, log *lg.Logger, config *conf.Config, pinmsg_log *tgbot.PinMsgFile) *scheduler.Job {
	curTime := getCurTime(config.Poll.Poll_start_time, "N", log)

	wednesdayJob, err := scheduler.Every().Day().At(curTime).Run(func() {
		// check day of week because sheduler.Every()DayOfWeek().Run() doesn't work
		// if time.Now().Weekday() == time.Wednesday {
		// 	lib.SendPoll(bot, &tb.Chat{ID: int64(config.Bot_secure.Chat_id)}, fact, log, config, pinmsg_log)
		// }
		tgbot.SendPoll(bot, &tb.Chat{ID: int64(config.Bot_secure.Chat_id)}, fact, log, config, pinmsg_log)
	})
	if err != nil {
		log.Fatal(err, "Ошибка при планировании задачи на среду:")
	}

	log.Info("Запланирована задача на среду")

	return wednesdayJob
}

func (s *Sheduler) ThursdayJob(bot *tb.Bot, fact string, log *lg.Logger, config *conf.Config, pinmsg_log *tgbot.PinMsgFile) *scheduler.Job {
	curTime := getCurTime(config.Poll.Poll_end_time, "Y", log)

	thursdayJob, err := scheduler.Every().Day().At(curTime).Run(func() {
		// check day of week because sheduler.Every()DayOfWeek().Run() doesn't work
		// if time.Now().Weekday() == time.Thursday {
		// 	lib.UnpinMsg(bot, &tb.Chat{ID: int64(config.Bot_secure.Chat_id)}, log, config, pinmsg_log)
		// }
		tgbot.UnpinMsg(bot, &tb.Chat{ID: int64(config.Bot_secure.Chat_id)}, log, config, pinmsg_log)
	})
	if err != nil {
		log.Fatal(err, "Ошибка при планировании задачи на четверг:")
	}

	log.Info("Запланирована задача на четверг")

	return thursdayJob
}

func getCurTime(configTime string, endFlg string, log *lg.Logger) string {
	// Convert poll start time to Time
	pollStartTimeT, err := time.Parse("15:04", configTime)
	if err != nil {
		log.Fatal(err, "Ошибка при конвертации времени начала опроса:")
		return ""
	}

	// Get current time in format HH24:MI
	curTime := time.Now().Format("15:04")

	// Convert curTime to Time
	curTimeT, err := time.Parse("15:04", curTime)
	if err != nil {
		log.Fatal(err, "Ошибка при конвертации текущего времени:")
		return ""
	}

	if curTimeT.After(pollStartTimeT) {
		// Add n minutes
		if endFlg == "Y" {
			curTimeT = curTimeT.Add(time.Minute * 2)
		} else {
			curTimeT = curTimeT.Add(time.Minute * 1)
		}
	} else {
		curTimeT = pollStartTimeT
	}

	// Convert curTimeT to string in format HH24:MI
	curTime = curTimeT.Format("15:04")

	return curTime
}

package lib

import (
	"fmt"
	"time"

	"math/rand"

	conf "app.go/app/config"
	lg "app.go/app/lib/logger"
	tb "gopkg.in/telebot.v3"
)

func CreateBot(fact string, log *lg.Logger, config *conf.Config) *tb.Bot {
	// Get bot token
	bot_token := DecryptBotToken(config.Bot_secure.Bot_token, config.Bot_secure.Bot_hash, log)

	// Set preferences for bot
	pref := tb.Settings{
		Token: bot_token,
		Poller: &tb.LongPoller{
			// Time for bot reading messages
			Timeout:        time.Duration(config.Bot_secure.Upd_time) * time.Second,
			AllowedUpdates: []string{"message", "edited_message", "channel_post", "edited_channel_post"}},
	}

	// Create bot
	bot, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err, "Ошибка создания бота:")
	}

	startReadingTime := time.Now()

	// Add command's handlers for debug
	bot.Handle("/sendpoll", func(c tb.Context) error {
		// Check all message in queue, their time and if it is more than now then run logic!
		if c.Message().Time().After(startReadingTime) {
			sendPoll(bot, c.Chat(), fact, log, config)
			// Save poll msg ID
			// c.Set("poll_msg_id", msg.ID)
		}
		return nil
	})

	bot.Handle("/unpinpoll", func(c tb.Context) error {
		// Check all message in queue, their time and if it is more than now then run logic!
		if c.Message().Time().After(startReadingTime) {
			unpinMsg(bot, c.Chat(), log)
		}
		return nil
	})

	log.Info(fmt.Sprintf("Бот '%s' создан успешно", bot.Me.Username))

	return bot
}

func sendPoll(bot *tb.Bot, chat *tb.Chat, fact string, log *lg.Logger, config *conf.Config) {
	if config.Poll.Question == "" {
		log.Fatal(fmt.Errorf("отсутствует вопрос. Заполните его в конфигурационном файле"), "")
	}

	result_message := fmt.Sprintf("%s\n%s", fact, config.Poll.Question)

	// Get answers
	if len(config.Poll.AnswersYes) == 0 || len(config.Poll.AnswersNo) == 0 {
		log.Fatal(fmt.Errorf("отсутствуют варианты ответов. Заполните их в конфигурационном файле"), "")
	}
	rnd_src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rnd_src)
	rnd_num := r.Intn(len(config.Poll.AnswersYes))
	answers := []string{config.Poll.AnswersYes[rnd_num], config.Poll.AnswersNo[rnd_num]}
	// Set answer's options for poll
	poll_options := []tb.PollOption{
		{Text: answers[0], VoterCount: 0},
		{Text: answers[1], VoterCount: 0},
	}

	// Create poll message
	poll_msg := &tb.Poll{
		Question:  result_message,
		Options:   poll_options,
		Type:      tb.PollRegular,
		Anonymous: false,
	}

	// Send poll
	poll_message, err := bot.Send(
		// Send's what's chat
		// &tb.Chat{ID: int64(chat_id)},
		chat,
		// Send's poll message
		poll_msg,
	)
	if err != nil {
		log.Fatal(err, "Ошибка отправки опроса:")
		// return nil
	}

	log.Info(fmt.Sprintf("Опрос успешно отправлен c ID = %v", poll_message.ID))

	// Pin created poll message
	pinMsg(bot, chat, log, poll_message)

	// return poll_message
}

func pinMsg(bot *tb.Bot, c *tb.Chat, log *lg.Logger, new_poll *tb.Message) {
	chat, err := bot.ChatByID(c.ID)
	if err != nil {
		log.Fatal(err, "Ошибка получения чата:")
		return
	}

	// Get pin messages
	pinnedMessage := chat.PinnedMessage

	if pinnedMessage != nil && pinnedMessage.Poll != nil && pinnedMessage.ID == new_poll.ID {
		log.Info(fmt.Sprintf("Опрос уже закреплен с ID = %v", new_poll.ID))
		return
	} else {
		// Pin poll
		err := bot.Pin(new_poll, tb.Protected)
		if err != nil {
			log.Fatal(err, fmt.Sprintf("Ошибка закрепления опроса с ID = %v", new_poll.ID))
		}
		log.Info(fmt.Sprintf("Опрос закреплен с ID = %v", new_poll.ID))
	}
}

func unpinMsg(bot *tb.Bot, c *tb.Chat, log *lg.Logger) {
	chat, err := bot.ChatByID(c.ID)
	if err != nil {
		log.Fatal(err, "Ошибка получения чата:")
		return
	}

	// Get pin messages
	pinnedMessage := chat.PinnedMessage

	if pinnedMessage != nil && pinnedMessage.Poll != nil {
		// Close/stop poll
		poll_msg, err := bot.StopPoll(pinnedMessage, tb.Silent)
		if err != nil {
			log.Fatal(err, "Ошибка закрытия опроса:")
		}
		log.Info(fmt.Sprintf("Опрос закрыт c ID = %v", pinnedMessage.ID))

		// Unpin poll
		err = bot.Unpin(chat, pinnedMessage.ID)
		if err != nil {
			log.Fatal(err, "Ошибка открепления опроса:")
		}
		log.Info(fmt.Sprintf("Опрос откреплен c ID = %v", pinnedMessage.ID))

		// Send result reply message
		// Get poll's results
		yesNum := poll_msg.Options[0].VoterCount
		noNum := poll_msg.Options[1].VoterCount

		// Check results and send result reply message
		if yesNum > noNum {
			log.Info(fmt.Sprintf("Положительных ответов: %d", yesNum))

			_, err = bot.Send(chat, "Ква, по результатам опроса встреча чуваков актуальна!", &tb.SendOptions{ReplyTo: pinnedMessage})
			if err != nil {
				log.Fatal(err, "Ошибка отправки результирующего сообщения по опросу:")
			}
		} else if noNum > yesNum {
			log.Info(fmt.Sprintf("Отрицательных ответов: %d", noNum))

			_, err = bot.Send(chat, "Ква, по результатам опроса встреча чуваков НЕ актуальна...", &tb.SendOptions{ReplyTo: pinnedMessage})
			if err != nil {
				log.Fatal(err, "Ошибка отправки результирующего сообщения по опросу:")
			}
		} else {
			log.Info("Нет очевидного результата")
		}
	}
}

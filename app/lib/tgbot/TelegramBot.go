package tgbot

import (
	"fmt"
	"time"

	"math/rand"

	conf "app.go/app/config"
	lib "app.go/app/lib"
	lg "app.go/app/lib/logger"
	tb "gopkg.in/telebot.v3"
)

func stopBotAfterExec(bot *tb.Bot, log *lg.Logger, config *conf.Config) {
	if config.Bot_secure.Exit_after_exec {
		log.Info("Бот остановлен")
		bot.Stop()
	}
}

func checkExistingChat(bot *tb.Bot, c *tb.Chat, log *lg.Logger) *tb.Chat {
	// Check existing chat
	chat, err := bot.ChatByID(c.ID)
	if err != nil {
		log.Fatal(err, "Ошибка получения чата:")
		return nil
	}

	return chat
}

func CreateBot(fact string, log *lg.Logger, config *conf.Config, pinmsg_log *PinMsgFile) *tb.Bot {
	// Get chatID
	// chatID := int64(config.Bot_secure.Chat_id)

	// Get bot token
	bot_token := lib.DecryptBotToken(config.Bot_secure.Bot_token, config.Bot_secure.Bot_hash, log)

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

	// Add command's handlers for debug
	/*startReadingTime := time.Now()
	bot.Handle("/sendpoll", func(c tb.Context) error {
		// Check all message in queue, their time and if it is more than now then run logic!
		if c.Message().Time().After(startReadingTime) {
			SendPoll(bot, c.Chat(), fact, log, config, pinmsg_log)
			// Save poll msg ID
			// c.Set("poll_msg_id", msg.ID)
		}
		return nil
	})

	bot.Handle("/unpinpoll", func(c tb.Context) error {
		// Check all message in queue, their time and if it is more than now then run logic!
		if c.Message().Time().After(startReadingTime) {
			UnpinMsg(bot, c.Chat(), log, config, pinmsg_log)
		}
		return nil
	})*/

	log.Info(fmt.Sprintf("Бот '%s' создан успешно", bot.Me.Username))

	return bot
}

func SendPoll(bot *tb.Bot, c *tb.Chat, fact string, log *lg.Logger, config *conf.Config, pinmsg_log *PinMsgFile) {
	// Check config Question
	if config.Poll.Question == "" {
		log.Fatal(fmt.Errorf("отсутствует вопрос. Заполните его в конфигурационном файле"), "")
	}

	// Check existing chat
	chat := checkExistingChat(bot, c, log)

	result_message := fmt.Sprintf("%s\n%s", fact, config.Poll.Question)

	// Get config answers
	if len(config.Poll.AnswersYes) == 0 || len(config.Poll.AnswersNo) == 0 {
		log.Fatal(fmt.Errorf("отсутствуют варианты ответов. Заполните их в конфигурационном файле"), "")
	}
	rnd_src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rnd_src)
	rnd_num := r.Intn(len(config.Poll.AnswersYes))
	answers := []string{config.Poll.AnswersYes[rnd_num], config.Poll.AnswersNo[rnd_num]}
	// Set answer options for poll
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
		// Send to what's chat
		chat,
		// Send message
		poll_msg,
	)
	if err != nil {
		log.Fatal(err, "Ошибка отправки опроса:")
		return
	}

	log.Info(fmt.Sprintf("Опрос успешно отправлен c ID = %v", poll_message.ID))

	// Pin created poll message
	pinMsg(bot, chat, log, poll_message)

	// Save poll msg ID
	pinmsg_log.Write(poll_message.ID, log)

	// Stop bot after command
	stopBotAfterExec(bot, log, config)
}

func pinMsg(bot *tb.Bot, c *tb.Chat, log *lg.Logger, new_poll *tb.Message) {
	// Check existing chat
	chat := checkExistingChat(bot, c, log)

	// Get pin message
	pinnedMessage := chat.PinnedMessage

	// Check if created poll already pinned
	if pinnedMessage != nil && pinnedMessage.Poll != nil && pinnedMessage.ID == new_poll.ID {
		log.Info(fmt.Sprintf("Опрос уже закреплен с ID = %v", new_poll.ID))
		return
	} else {
		// Pin poll
		err := bot.Pin(new_poll, tb.Protected)
		if err != nil {
			log.Fatal(err, fmt.Sprintf("Ошибка закрепления опроса с ID = %v", new_poll.ID))
			return
		}

		log.Info(fmt.Sprintf("Опрос закреплен с ID = %v", new_poll.ID))
	}
}

func UnpinMsg(bot *tb.Bot, c *tb.Chat, log *lg.Logger, config *conf.Config, pinmsg_log *PinMsgFile) {
	// Check existing chat
	chat := checkExistingChat(bot, c, log)

	// Get pin messages
	pinnedMessage := chat.PinnedMessage
	// Get last poll msg ID
	logPinnedMessageID := pinmsg_log.GetLastPollMsgID()
	// pinnedMessage := &tb.Message{ID: logPinnedMessageID, Chat: chat}
	// Find message by ID

	log.Info(fmt.Sprintf("Найдено сообщение: %s", pinnedMessage.Poll.Question))

	// Check if created poll already pinned

	if pinnedMessage.Poll != nil && pinnedMessage.ID == logPinnedMessageID {
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
		// Get poll results
		yesNum := poll_msg.Options[0].VoterCount
		noNum := poll_msg.Options[1].VoterCount
		var result_msg string

		// Check results and send result reply message
		if yesNum > noNum && yesNum > 2 {
			log.Info(fmt.Sprintf("Положительных ответов: %d", yesNum))

			result_msg = "Ква, по результатам опроса встреча чуваков актуальна!"
			_, err = bot.Send(chat, result_msg, &tb.SendOptions{ReplyTo: pinnedMessage})
			if err != nil {
				log.Fatal(err, "Ошибка отправки результирующего сообщения по опросу:")
			}

			log.Info(fmt.Sprintf("Результирующее сообщение отправлено: %s", result_msg))
		} else if noNum > yesNum {
			log.Info(fmt.Sprintf("Отрицательных ответов: %d", noNum))

			result_msg = "Ква, по результатам опроса встреча чуваков НЕ актуальна..."
			_, err = bot.Send(chat, result_msg, &tb.SendOptions{ReplyTo: pinnedMessage})
			if err != nil {
				log.Fatal(err, "Ошибка отправки результирующего сообщения по опросу:")
			}

			log.Info(fmt.Sprintf("Результирующее сообщение отправлено: %s", result_msg))
		} else {
			log.Info("Нет очевидного результата")

			result_msg = "Ква, не поня-я-тно. Решайте сами, чуваки"
			_, err = bot.Send(chat, result_msg, &tb.SendOptions{ReplyTo: pinnedMessage})
			if err != nil {
				log.Fatal(err, "Ошибка отправки результирующего сообщения по опросу:")
			}

			log.Info(fmt.Sprintf("Результирующее сообщение отправлено: %s", result_msg))
		}
	} else {
		log.Info("Опрос ранее не был закреплен")
	}

	pinmsg_log.ClearFile()

	// Stop bot after command
	stopBotAfterExec(bot, log, config)
}

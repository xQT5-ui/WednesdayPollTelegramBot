package lib

import (
	"fmt"
	"time"

	conf "app.go/app/config"
	lg "app.go/app/lib/logger"
	"golang.org/x/exp/rand"
	tb "gopkg.in/telebot.v3"
)

func CreateBot(chat_id int, fact, question string, answersYes, answersNo []string, log *lg.Logger, config *conf.Config) *tb.Bot {
	// Get bot token
	bot_token := DecryptBotToken(config.Bot_secure.Bot_token, config.Bot_secure.Bot_hash, log)

	// Set preferences for bot
	pref := tb.Settings{
		Token: bot_token,
		Poller: &tb.LongPoller{
			// Time for bot reading messages
			Timeout:        10 * time.Second,
			AllowedUpdates: []string{"message", "edited_message", "channel_post", "edited_channel_post"}},
	}

	// Create bot
	bot, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err, "Ошибка создания бота:")
	}

	startReadingTime := time.Now()
	var msg *tb.Message

	// Add command's handlers for debug
	bot.Handle("/sendpoll", func(c tb.Context) error {
		// Check all message in queue, their time and if it is more than now then run logic!
		if c.Message().Time().After(startReadingTime) {
			msg = sendPoll(bot, c.Chat(), fact, question, answersYes, answersNo, log)
			// Save poll msg ID
			// c.Set("poll_msg_id", msg.ID)
		}
		return nil
	})

	bot.Handle("/pinpoll", func(c tb.Context) error {
		// Check all message in queue, their time and if it is more than now then run logic!
		if c.Message().Time().After(startReadingTime) {
			pinMsg(bot, c.Chat(), log, msg)
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

func sendPoll(bot *tb.Bot, chat *tb.Chat, fact, question string, answersYes, answersNo []string, log *lg.Logger) *tb.Message {
	if question == "" {
		log.Fatal(fmt.Errorf("отсутствует вопрос. Заполните его в конфигурационном файле"), "")
	}

	result_message := fmt.Sprintf("%s\n%s", fact, question)

	// Get answers
	if len(answersYes) == 0 || len(answersNo) == 0 {
		log.Fatal(fmt.Errorf("отсутствуют варианты ответов. Заполните их в конфигурационном файле"), "")
	}
	rand.Seed(34)
	rnd_num := rand.Intn(len(answersYes))
	answers := []string{answersYes[rnd_num], answersNo[rnd_num]}
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
		return nil
	}

	log.Info(fmt.Sprintf("Опрос успешно отправлен:\n%v", poll_message.ID))

	return poll_message
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
		log.Info(fmt.Sprintf("Опрос уже закреплен с ID = '%v'", new_poll.ID))
		return
	} else {
		// Pin poll
		err := bot.Pin(new_poll, tb.Silent)
		if err != nil {
			log.Fatal(err, fmt.Sprintf("Ошибка закрепления опроса с ID = '%v'", new_poll.ID))
		}
		log.Info(fmt.Sprintf("Опрос закреплен ID = '%v'", new_poll.ID))
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
		// Unpin poll
		err := bot.Unpin(chat, pinnedMessage.ID)
		if err != nil {
			log.Fatal(err, "Ошибка открепления опроса:")
		}
		log.Info(fmt.Sprintf("Опрос откреплен c ID = '%v'", pinnedMessage.ID))
	}
}

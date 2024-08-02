package lib

import (
	"fmt"
	"time"

	conf "app.go/app/config"
	lg "app.go/app/lib/logger"
	"golang.org/x/exp/rand"
	tb "gopkg.in/tucnak/telebot.v2"
)

func CreateBot(chat_id int, fact, question string, answersYes, answersNo []string, log *lg.Logger, config *conf.Config) *tb.Bot {
	// Get bot token
	bot_token := DecryptBotToken(config.Bot_secure.Bot_token, config.Bot_secure.Bot_hash, log)

	bot, err := tb.NewBot(tb.Settings{
		Token: bot_token,
		// Time for bot reading messages
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err, "Ошибка создания бота:")
	}

	// Add command's handlers for debug
	bot.Handle("/sendpoll", func(m *tb.Message) {
		SendPoll(bot, chat_id, fact, question, answersYes, answersNo, log)
	})

	log.Info(fmt.Sprintf("Бот '%s' создан успешно", bot.Me.Username))

	return bot
}

func SendPoll(bot *tb.Bot, chat_id int, fact, question string, answersYes, answersNo []string, log *lg.Logger) {
	if question == "" {
		log.Fatal(fmt.Errorf("отсутствует вопрос. Заполните его в конфигурационном файле"), "")
	}
	result_message := fmt.Sprintf("%s\n%s", fact, question)

	// Get answers
	if len(answersYes) == 0 || len(answersNo) == 0 {
		log.Fatal(fmt.Errorf("отсутствуют варианты ответов. Заполните их в конфигурационном файле"), "")
	}
	rnd_num := rand.Intn(len(answersYes))
	answers := []string{answersYes[rnd_num], answersNo[rnd_num]}
	// Set answer's options for poll
	poll_options := []tb.PollOption{
		{Text: answers[0], VoterCount: 0},
		{Text: answers[1], VoterCount: 0},
	}

	// Send poll
	poll_message, err := bot.Send(
		// Send to what's chat
		&tb.Chat{ID: int64(chat_id)},
		// Send poll message
		&tb.Poll{
			Question:  result_message,
			Options:   poll_options,
			Type:      tb.PollRegular,
			Anonymous: false,
		},
	)
	if err != nil {
		log.Fatal(err, "Ошибка отправки опроса:\n%v")
		// return nil
	}

	log.Info(fmt.Sprintf("Опрос успешно отправлен:\n%v", poll_message))

	// return poll_message
}

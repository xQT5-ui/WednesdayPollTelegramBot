package main

import (
	"app.go/app/config"
	"app.go/app/lib"
	"app.go/app/lib/logger"
)

func main() {
	//Get logger
	log := logger.NewLogger(nil)
	log.Info("Запуск бота")

	// Get config
	config := config.LoadConfig(log)

	// Get website data
	fact := lib.DataFromWebsite(config.Url, log)

	// Create bot
	bot := lib.CreateBot(config.Bot_secure.Chat_id, fact, config.Poll.Question, config.Poll.AnswersYes, config.Poll.AnswersNo, log, &config)
	log.Info("Бот запущен")

	// Start bot in infinite loop
	bot.Start()
}

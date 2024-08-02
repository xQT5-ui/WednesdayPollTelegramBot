package main

import (
	"app.go/app/config"
	"app.go/app/lib"
)

func main() {
	// Get config
	config := config.LoadConfig()

	// Get bot_token
	bot_token := lib.GetDecryptBotToken(config.Bot_secure.Bot_token, config.Bot_secure.Bot_hash)

	// Get website data
	fact := lib.GetDataFromWebsite(config.Url)

	// Create bot
	bot := lib.CreateBot(bot_token, config.Bot_secure.Chat_id, fact, config.Poll.Question, config.Poll.AnswersYes, config.Poll.AnswersNo)

	// Start bot in infinite loop
	bot.Start()
}

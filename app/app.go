package main

import (
	"os"
	"os/signal"
	"syscall"

	"app.go/app/config"
	"app.go/app/lib"
	"app.go/app/lib/logger"
)

func main() {
	//Get logger
	log := logger.NewLogger(nil)
	// Close open file
	defer log.Close()
	log.Info("Запуск бота")

	// Create gorutine for logger
	c_sys := make(chan os.Signal, 1)
	// Set system signal to channel
	signal.Notify(c_sys, os.Interrupt, syscall.SIGTERM)
	go func() {
		// Waiting signal from channel (from signal.Notify)
		<-c_sys
		log.Close()
		os.Exit(0)
	}()

	// Get config
	config := config.LoadConfig(log)

	// Get website data
	fact := lib.DataFromWebsite(config.Url, log)

	// Create bot
	bot := lib.CreateBot(fact, log, &config)
	log.Info("Бот запущен")

	// Start bot in infinite loop
	bot.Start()
}

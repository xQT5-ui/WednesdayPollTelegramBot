package main

import (
	"os"
	"os/signal"
	"syscall"

	"app.go/app/config"
	"app.go/app/lib"
	"app.go/app/lib/logger"
	"app.go/app/lib/sheduler"
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

	// Get config
	config := config.LoadConfig(log)

	// Get website data
	fact := lib.DataFromWebsite(config.Url, log)
	// fact := "В состав архипелага Филиппины входит 7107 островов."

	// Create bot
	bot := lib.CreateBot(fact, log, &config)

	// Create sheduler and jobs
	sheduler := sheduler.NewScheduler()
	log.Info("Планировщик создан")
	wednesdayJob := sheduler.WednesdayJob(bot, fact, log, &config)
	thursdayJob := sheduler.ThursdayJob(bot, fact, log, &config)

	// FOR TEST: Wait for system signal
	/*go func() {
		time.Sleep(5 * time.Second)
		log.Info("Прошло 10 секунд. Отправка сигнала остановки программы...")
		c_sys <- os.Interrupt
	}()*/

	go func() {
		// Waiting signal from channel (from signal.Notify)
		<-c_sys
		// Stop jobs
		wednesdayJob.Quit <- true
		thursdayJob.Quit <- true
		log.Info("Задачи остановлены")
		// Close log file
		log.Close()
		// Stop bot
		log.Info("Бот остановлен")
		bot.Stop()
		// Exit program
		os.Exit(0)
	}()

	// Start bot in infinite loop
	log.Info("Бот запущен")
	bot.Start()
}

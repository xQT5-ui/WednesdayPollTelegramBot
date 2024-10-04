package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"app.go/app/lib/emptyfile"
)

type Logger struct {
	Log    *log.Logger
	File   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
}

func NewLogger(l *log.Logger) *Logger {
	if l == nil {
		l = log.Default()
	}

	file := emptyfile.ReCreateFile("log", "WednesdayPollTgBot.log")

	return &Logger{
		Log:  l,
		File: file,
		// Create new buffered writer for file: first program wirte text in buffer, then write in file for safety and minimize RAM
		writer: bufio.NewWriter(file),
	}
}

func (l *Logger) Info(message string) {
	// Lock for thread safety
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt_time := time.Now().Format("2006-01-02 15:04:05")
	log_msg := fmt.Sprintf("[%v] INFO: %s", fmt_time, message)

	l.Log.Println(log_msg)

	// Set message into file
	_, err := l.writer.WriteString(log_msg + "\n")
	if err != nil {
		l.Log.Fatalf("FATAL: %s\n%v", "Ошибка записи в файл:", err)
	}

	// Set text into file from buffer
	l.writer.Flush()
}

func (l *Logger) Fatal(err error, message string) {
	// Lock for thread safety
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt_time := time.Now().Format("2006-01-02 15:04:05")
	log_msg := fmt.Sprintf("[%v] FATAL: %s\n%v", fmt_time, message, err)

	l.Log.Fatalf(log_msg)

	// Set message into file
	_, err = l.writer.WriteString(log_msg + "\n")
	if err != nil {
		l.Log.Fatalf("FATAL: %s\n%v", "Ошибка записи в файл:", err)
	}

	// Set text into file from buffer
	l.writer.Flush()
}

func (l *Logger) Close() {
	// Lock for thread safety
	l.mu.Lock()
	defer l.mu.Unlock()

	// Clear buffer
	if l.writer != nil {
		l.writer.Flush()
	}
	// Close file
	if l.File != nil {
		l.File.Close()
	}
}

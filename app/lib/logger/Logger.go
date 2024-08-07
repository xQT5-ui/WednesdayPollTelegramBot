package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
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

	// Get executable path
	var file *os.File

	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err, "Ошибка получения пути исполняемого файла:")
		return nil
	}

	// Get executable dir
	execDir := filepath.Dir(execPath)

	// Set config path
	filename := filepath.Join(execDir, "log", "WednesdayPollTgBot.log")

	// Create file
	// Delete previous file if it exists
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		err = os.Remove(filename)
		if err != nil {
			log.Fatal(err, fmt.Sprintf("Не получилось удалить файл лога '%s':", filename))
			return nil
		}
	}
	file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err, fmt.Sprintf("Не получилось создать файл лога '%s':", filename))
		return nil
	}

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

	l.Log.Printf(log_msg)

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

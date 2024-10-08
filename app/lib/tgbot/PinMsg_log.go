package tgbot

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"app.go/app/lib/emptyfile"
	lg "app.go/app/lib/logger"
)

type PinMsgFile struct {
	File   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
}

func NewPinMsgFile(log *lg.Logger) *PinMsgFile {
	file := emptyfile.ReCreateFile("log", "WednesdayPollPinMsg.txt")

	return &PinMsgFile{
		File: file,
		// Create new buffered writer for file: first program wirte text in buffer, then write in file for safety and minimize RAM
		writer: bufio.NewWriter(file),
	}
}

func (l *PinMsgFile) Write(msg int, log *lg.Logger) {
	// Lock for thread safety
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt_time := time.Now().Format("2006-01-02")
	log_msg := fmt.Sprintf("[%v] PinMsg ID: %v", fmt_time, msg)

	log.Info(log_msg)

	// Set message into file
	_, err := l.writer.WriteString(log_msg + "\n")
	if err != nil {
		log.Fatal(err, "Ошибка записи в файл:")
	}

	// Set text into file from buffer
	l.writer.Flush()
}

func (l *PinMsgFile) ClearFile() {
	// Lock for thread safety
	l.mu.Lock()
	defer l.mu.Unlock()

	// Clear buffer
	if l.writer != nil {
		l.writer.Flush()
	}
	// Clear file
	if l.File != nil {
		l.File.Truncate(0)
		l.File.Seek(0, 0)
	}
}

func (l *PinMsgFile) Close() {
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

func (l *PinMsgFile) GetLastPollMsgID() int {
	// Lock for thread safety
	l.mu.Lock()
	defer l.mu.Unlock()

	// Get last line from file
	file, err := os.Open(l.File.Name())
	if err != nil {
		log.Fatal(err, "Не получилось открыть файл лога закреплённого сообщения:")
		return 0
	}
	defer file.Close()

	// Read last line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lastLine string

	for scanner.Scan() {
		lastLine = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err, "Не получилось прочитать файл лога закреплённого сообщения:")
		return 0
	}

	// Parse last line throw regexp
	re := regexp.MustCompile(`ID: (\d+)`)
	matches := re.FindStringSubmatch(lastLine)
	if len(matches) != 2 {
		log.Fatal(fmt.Errorf("длина закреплённого ID закрепленного сообщения меньше 2 элементов: %s", lastLine), "Не удалось распарсить ID закрепленного сообщения:")
	}
	id, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatal(err, "Не удалось конвертировать ID закрепленного сообщения:")
		return 0
	}

	return id
}

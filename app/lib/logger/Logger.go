package logger

import "log"

type Logger struct {
	Log *log.Logger
}

func NewLogger(l *log.Logger) *Logger {
	if l == nil {
		l = log.Default()
	}

	return &Logger{
		Log: l,
	}
}

func (l *Logger) Info(message string) {
	l.Log.Printf("INFO: %v", message)
}

func (l *Logger) Fatal(err error, message string) {
	l.Log.Fatalf("FATAL: %s\n%v", message, err)
}

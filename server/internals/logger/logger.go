package logger

import (
	"log"
)

type Logger struct {
	info *log.Logger
	error *log.Logger
}

func New(info, error *log.Logger) *Logger {
	return &Logger{
		info: info,
		error: error,
	}
}

func (l *Logger) Info(msg string, fields map[string]any) {
	l.info.Println(msg,fields)
}

func (l *Logger) Error(msg string, fields map[string]any) {
	l.info.Println(msg,fields)
}
package logger

import (
	"fmt"
	"github.com/bonjourmalware/kuping/internal/router"
	"os"
)

type Logger struct {
	File *os.File
}

func NewLogger(filepath string) *Logger {
	logfile, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to open log file [%s]", filepath))
		os.Exit(1)
	}

	l := &Logger{
		File: logfile,
	}

	return l
}

func (l *Logger) Write(rawData router.Event) error {
	data, err := rawData.String()
	if err != nil {
		return err
	}

	_, err = l.File.WriteString(data + "\n")
	if err != nil {
		return err
	}

	return nil
}

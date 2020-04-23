package logger

import (
	"fmt"
	"gitlab.com/Alvoras/kuping/internal/router"
	"time"
)

func Start(loggerChan chan router.Event, logfilePath string) {
	go receiveEvents(loggerChan, logfilePath)
}

func receiveEvents(loggerChan chan router.Event, logfilePath string) {
	logWriter := NewLogger(logfilePath)

	for {
		select {
		case ev := <-loggerChan:
			err := logWriter.Write(ev)
			if err != nil {
				fmt.Println("failed to stringify JSON payload while writing to log file", time.Now().Format(time.RFC3339))
				continue
			}
		}
	}
}

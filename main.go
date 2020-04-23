package main

import (
	"fmt"
	"gitlab.com/Alvoras/kuping/internal/filter"
	"gitlab.com/Alvoras/kuping/internal/logger"
	"os"

	"gitlab.com/Alvoras/kuping/internal/config"
	"gitlab.com/Alvoras/kuping/internal/router"
)

func init() {
	config.Cfg.Load()
}

func main() {
	quitErrChan := make(chan error)
	logChan := make(chan router.Event)
	eventChan := make(chan router.Event)

	logger.Start(logChan, config.Cfg.Logfile)
	filter.Start(logChan, eventChan)
	router.StartServers(eventChan, quitErrChan)
	fmt.Println("OK")

	select {
	case err := <-quitErrChan:
		fmt.Println(err)
		os.Exit(1)
	}
}
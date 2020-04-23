package router

import (
	"crypto/tls"
	"fmt"
	"gitlab.com/Alvoras/kuping/internal/config"
	"net/http"
)

func StartServers(eventChan chan Event, quitErrChan chan error) {
	if config.Cfg.HTTP.Enabled {
		// For every HTTP port submitted, start an HTTP listener on that port
		fmt.Println("Starting HTTP servers...")
		for port := 1; port <= 65535; port++ {
			if config.Cfg.HTTP.EnabledPorts.Contains(uint64(port)){
				go startHTTP(port, quitErrChan, eventChan)
			}
		}
	}

	if config.Cfg.HTTPS.Enabled {
		fmt.Println("Starting HTTPS servers...")
		for port := 1; port <= 65535; port++ {
			if config.Cfg.HTTPS.EnabledPorts.Contains(uint64(port)){
				go startHTTPS(port, quitErrChan, eventChan)
			}
		}
	}
}

func startHTTP(port int, quitErrChan chan error, eventChan chan Event) {
	HTTPRouter := http.NewServeMux()

	handler := http.HandlerFunc(Index)
	HTTPRouter.Handle("/", logHandler(handler, eventChan))

	fmt.Println("Started HTTP server on port :", port)
	quitErrChan <- http.ListenAndServe(fmt.Sprintf(":%d", port), HTTPRouter)
}


func startHTTPS(port int, quitErrChan chan error, eventChan chan Event) {
	HTTPSRouter := http.NewServeMux()
	handler := http.HandlerFunc(IndexHTTPS)
	HTTPSRouter.Handle("/", logHandler(handler, eventChan))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      HTTPSRouter,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	fmt.Println("Started HTTPS server on port :", port)
	quitErrChan <- srv.ListenAndServeTLS("server.crt", "server.key")
}

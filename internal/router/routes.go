package router

import (
	"net/http"

	"gitlab.com/Alvoras/kuping/internal/config"
)

func Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Server", config.Cfg.ServerHeader)
}

func IndexHTTPS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Server", config.Cfg.ServerHeader)
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
}

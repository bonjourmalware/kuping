package router

import (
	"net/http"
)

func Index(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Server", "Apache")
}

func IndexHTTPS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Server", "Apache")
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
}

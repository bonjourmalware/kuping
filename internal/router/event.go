package router

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Timestamp  string `json:"timestamp"`
	Verb       string `json:"verb"`
	Proto      string `json:"proto"`
	RequestURI string `json:"URI"`
	RemoteAddr string `json:"remote_address"`
	SourceIP   string `json:"src_ip"`
	SourcePort uint64 `json:"src_port"`
	DestPort   uint64 `json:"dst_port"`
	DestHost   string `json:"dst_host"`
	Host       string `json:"host"`
	StatusCode int `json:"status_code"`
	Headers map[string]string `json:"headers"`
	Errors  []string          `json:"errors"`
	Params  string            `json:"params"`
	IsTLS   bool              `json:"is_tls"`
}

func newEvent(r *http.Request) Event {
	var srcIP string
	var srcPortStr string
	var srcPort uint64
	var dstPortStr string
	var dstPort uint64
	var dstHost string
	headers := make(map[string]string)
	var errs []string
	var params []byte
	var isTLS bool

	for header := range r.Header {
		headers[header] = r.Header.Get(header)
	}

	srcIP = r.RemoteAddr

	if strings.Contains(r.RemoteAddr, ":") {
		rAddrSplit := strings.Split(r.RemoteAddr, ":")
		srcIP, srcPortStr = rAddrSplit[0], rAddrSplit[1]
		srcPort, _ = strconv.ParseUint(srcPortStr, 10, 16)
	}

	params, err := GetBodyPayload(r)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if strings.Contains(r.Host, ":") {
		dstHost, dstPortStr, err = net.SplitHostPort(r.Host)
		if err == nil {
			dstPort, _ = strconv.ParseUint(dstPortStr, 10, 16)
		}
	} else {
		if r.TLS != nil {
			// Default HTTPS port
			dstPort = 443
		} else {
			// Default HTTP port
			dstPort = 80
		}
	}

	if r.TLS != nil {
		isTLS = true
	} else {
		isTLS = false
	}

	data := Event{
		Timestamp:  time.Now().Format(time.RFC3339),
		Verb:       r.Method,
		Proto:      r.Proto,
		RequestURI: r.URL.RequestURI(),
		RemoteAddr: r.RemoteAddr,
		SourceIP:   srcIP,
		SourcePort: srcPort,
		DestPort:   dstPort,
		DestHost:   dstHost,
		Host:       r.Host,
		IsTLS:      isTLS,
		StatusCode: http.StatusOK, // Default, can be replaced if we implement dynamic response code
		Headers: headers,
		Errors:  errs,
		Params:  string(params),
	}

	return data
}

func (ev Event) String() (string, error) {
	data, err := json.Marshal(ev)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

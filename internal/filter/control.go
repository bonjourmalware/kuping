package filter

import (
	"gitlab.com/Alvoras/kuping/internal/config"
	"gitlab.com/Alvoras/kuping/internal/router"
)

func Start(loggerChan chan router.Event, eventChan chan router.Event) {
	go receiveEvents(loggerChan, eventChan)
}

func receiveEvents(loggerChan chan router.Event, eventChan chan router.Event) {
	for {
		select {
		case ev := <-eventChan:
			if ok, err := filterEvents(ev); err == nil && ok {
				loggerChan <- ev
			}
		}
	}
}

func filterEvents(ev router.Event) (bool, error) {
	// If no filtering list are enabled, log all
	if !config.Cfg.EnableWhitelist && !config.Cfg.EnableBlacklist{
		return true, nil
	}

	ok := false

	var wl config.IPRanges
	var bl config.IPRanges

	if ev.IsTLS {
		wl = config.Cfg.HTTPS.IPRules.WhitelistedIPs
		bl = config.Cfg.HTTPS.IPRules.BlacklistedIPs
	}else{
		wl = config.Cfg.HTTP.IPRules.WhitelistedIPs
		bl = config.Cfg.HTTP.IPRules.BlacklistedIPs
	}

	if config.Cfg.EnableWhitelist {
		for _, iprange := range wl {
			if iprange.ContainsIPString(ev.SourceIP) {
				ok = true

			}
		}
	}

	if config.Cfg.EnableBlacklist {
		for _, iprange := range bl {
			if iprange.ContainsIPString(ev.SourceIP) {
				ok = false

			}
		}
	}

	return ok, nil
}

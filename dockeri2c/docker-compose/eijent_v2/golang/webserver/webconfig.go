package webserver

import (
	"net/http"
)

type WebConfig struct {
	Pass    string
	Handler func(http.ResponseWriter, *http.Request)
}

func Config(cfg *SetupServer, wconfs []WebConfig) error {
	for _, wconf := range wconfs {
		cfg.Add(wconf.Pass, wconf.Handler)
	}
	return nil
}
